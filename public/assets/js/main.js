const W = 13;
const H = W;
const TOTAL_CELL = W * H;

const draw = (ctx) => {
    ctx.save();
    for (let r = 0; r < H; r++) {
        for (let c = 0; c < W; c++) {
            const idx = r * W + c;
            const deg = (360 * idx) / TOTAL_CELL;
            ctx.fillStyle = `hsl(${deg}deg, 100%, 50%)`;
            ctx.fillRect(c, r, 1, 1);
        }
    }
    ctx.restore();
};

const setup = () => {
    console.log("setup start");
    const canvas = document.querySelector("canvas");
    canvas.width = W;
    canvas.height = H;
    const ctx = canvas.getContext("2d");
    ctx.fillRect(0, 0, W, H);
    draw(ctx);
    const imgData = ctx.getImageData(0,0,W,H);
    console.log("setup end");
    return imgData;
};

const main = () => {
    console.log('main start');
    const imgData = setup();
    console.log('imgData', imgData);
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
            ws.send(imgData.data);
        }
    };

    ws.addEventListener('error', (e) => console.error('websocket error:', e));
    ws.addEventListener('open', (e) => console.log('websocket open:', e));
    ws.addEventListener('close', (e) => console.log('websocket close:', e));
    ws.addEventListener('message', handler_msg);
    console.log('main end');
};

main();
