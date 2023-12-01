const W = 13;
const H = W;
const TOTAL_CELL = W * H;

const fract = (x) => x - Math.floor(x);

const draw = (ctx, _t) => {
    // const t = _t * 0.001;
    const t = _t * 0.001;
    ctx.save();
    for (let r = 0; r < H; r++) {
        const offset_x = Math.floor(Math.sin(r/6+t) * 13);
        for (let c = 0; c < W-offset_x; c++) {
            const idx = r * W + c;
            const deg = 360 * fract(t + idx / TOTAL_CELL);
            ctx.fillStyle = `hsl(${deg}deg, 100%, 50%)`;
            ctx.fillRect(c, r, 1, 1);
        }
        for (let c = W-offset_x; c < W; c++) {
            ctx.fillStyle = 'black';
            ctx.fillRect(c, r, 1, 1);
        }
    }
    // ctx.fillStyle = 'black';
    // ctx.fillRect(0, 0, W, H);
    // ctx.fillStyle = '#0000ff';
    // ctx.font = "13px serif";
    // const text = "Happy Dussehra!!";
    // const text_len = text.length * 6;
    // ctx.fillText(text, 10 + fract(t) * -text_len, 10);
    ctx.restore();
};

const main = () => {
    console.log('main start');
    const ws = new WebSocket(`ws://${window.location.host}/ws`);
    ws.binaryType = "arraybuffer";
    const handler_msg = (e) => {
        console.log('websocket message:', e);
        const is_binary = e.data instanceof ArrayBuffer;
        console.log('is_binary', is_binary);
        if (is_binary) {
            const dec = new TextDecoder('utf-8');
            const s = dec.decode(e.data);
            console.log('decode binary message as utf-8 string:', s);
            ws.send('hello from websocket client');
        }
    };

    const setup = () => {
        console.log("setup start");
        const canvas = document.querySelector("canvas");
        canvas.width = W;
        canvas.height = H;
        const ctx = canvas.getContext("2d", { willReadFrequently: true });
        // ctx.rotate(Math.PI);
        // ctx.translate(0, 0);
        ctx.scale(-1, -1);
        ctx.translate(-W, -H);
        ctx.fillRect(0, 0, W, H);
        let last_t = 0;
        const step = (t) => {
            requestAnimationFrame(step);
            if (t - last_t < 1) return;
            last_t = t;
            draw(ctx, t);
            const imgData = ctx.getImageData(0, 0, W, H);
            // console.log('imgData:', imgData);
            ws.send(imgData.data);
        };
        requestAnimationFrame(step);
        console.log("setup end");
    };

    ws.addEventListener('error', (e) => console.error('websocket error:', e));
    ws.addEventListener('open', (e) => {
        console.log('websocket open:', e);
        setup();
    });
    ws.addEventListener('close', (e) => console.log('websocket close:', e));
    ws.addEventListener('message', handler_msg);


    console.log('main end');
};

main();
