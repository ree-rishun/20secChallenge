const cvs = document.querySelector("#draw-area");
const ctx = cvs.getContext("2d");

cvs.addEventListener("mousemove", event => {
  draw(event.layerX, event.layerY);
});
cvs.addEventListener("touchmove", event => {
  draw(event.layerX, event.layerY);
});

cvs.addEventListener("mousedown", () => {
  ctx.beginPath();
  isDrag = true;
});
cvs.addEventListener("mouseup", () => {
  ctx.closePath();
  isDrag = false;
});
cvs.addEventListener("touchstart", () => {
  ctx.beginPath();
  isDrag = true;
});
cvs.addEventListener("touchend", () => {
  ctx.closePath();
  isDrag = false;
});

const clearButton = document.querySelector("#clear-button");
clearButton.addEventListener("click", () => {
  ctx.fillStyle = 'rgb(255, 255,255)';
  ctx.fillRect(0, 0, cvs.width, cvs.height);
});

let isDrag = false;
function draw(x, y) {
  if (!isDrag) {
    return;
  }

  ctx.lineWidth = 5;
  ctx.lineTo(x, y);
  ctx.stroke();
}