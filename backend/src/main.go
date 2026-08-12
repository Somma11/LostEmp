package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		html := `<!doctype html>
<html lang="pt-BR">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>foda-se deu vontade</title>
<style>
  :root{
    --bg1:#0a0014; --bg2:#1f005c; --accent1:#ff4d6d; --accent2:#6a11cb; --accent3:#00d2ff; --gold:#ffd700;
  }
  *{box-sizing:border-box}
  html,body{height:100%;margin:0}
  body{
    position:relative;
    display:flex;align-items:center;justify-content:center;
    min-height:100vh;
    background:
      radial-gradient(circle at 15% 20%, rgba(255,77,109,0.25), transparent 35%),
      radial-gradient(circle at 85% 80%, rgba(0,210,255,0.25), transparent 35%),
      linear-gradient(135deg,var(--bg1),var(--bg2) 40%,#3d0066 70%,#0a0014);
    background-size:200% 200%;
    animation:bgshift 18s ease-in-out infinite;
    font-family:"Segoe UI",Inter,system-ui,Roboto,Helvetica,Arial,sans-serif;
    color:#fff;overflow:hidden;cursor:none;
  }
  @keyframes bgshift{
    0%,100%{background-position:0% 50%}
    50%{background-position:100% 50%}
  }

  canvas{position:fixed;inset:0;display:block}
  #stars{z-index:0;opacity:0.8}
  #matrix{z-index:1;opacity:0.05;mix-blend-mode:screen}
  #confetti{z-index:5;pointer-events:none}
  #cursor{z-index:9999;pointer-events:none}

  .stage{
    position:relative;z-index:3;text-align:center;
    padding:60px 50px;border-radius:28px;
    background:rgba(255,255,255,0.04);
    border:1px solid rgba(255,255,255,0.12);
    backdrop-filter:blur(14px) saturate(150%);
    box-shadow:0 20px 80px rgba(0,0,0,0.55), 0 0 120px rgba(255,77,109,0.08) inset;
    animation:float 6s ease-in-out infinite;
  }
  @keyframes float{
    0%,100%{transform:translateY(0px) rotate(0deg)}
    50%{transform:translateY(-10px) rotate(0.3deg)}
  }

  .badge{
    display:inline-block;margin-bottom:18px;padding:6px 16px;border-radius:999px;
    font-size:12px;letter-spacing:0.15em;text-transform:uppercase;
    background:linear-gradient(90deg,var(--accent1),var(--accent2),var(--accent3));
    background-size:300% 100%;animation:slide 4s linear infinite;
    box-shadow:0 6px 24px rgba(0,0,0,0.4);
  }
  @keyframes slide{0%{background-position:0% 0}100%{background-position:300% 0}}

  h1{
    font-size:clamp(34px,7vw,104px);
    margin:0 0 14px;letter-spacing:0.02em;text-transform:uppercase;
    font-weight:900;
    background:linear-gradient(90deg,#fff,var(--gold),var(--accent3),var(--accent1),#fff);
    background-size:400% auto;
    -webkit-background-clip:text;background-clip:text;color:transparent;
    animation:shine 6s linear infinite, pop 0.7s cubic-bezier(.2,1.6,.4,1);
    filter:drop-shadow(0 8px 30px rgba(0,0,0,0.6));
  }
  @keyframes shine{to{background-position:400% center}}
  @keyframes pop{from{transform:scale(0.6);opacity:0}to{transform:scale(1);opacity:1}}

  h1 span{display:inline-block;animation:bounce 2.4s ease-in-out infinite}
  h1 span:nth-child(odd){animation-delay:.15s}
  @keyframes bounce{0%,100%{transform:translateY(0)}50%{transform:translateY(-6px)}}

  .sub{
    margin:0 0 26px;font-size:clamp(14px,2vw,19px);opacity:0.85;
    font-style:italic;
  }
  .sub b{background:linear-gradient(90deg,var(--gold),var(--accent1));-webkit-background-clip:text;background-clip:text;color:transparent;font-style:normal}

  .counter{
    font-variant-numeric:tabular-nums;
    font-size:clamp(28px,4vw,46px);font-weight:800;margin:0 0 26px;
    text-shadow:0 0 30px rgba(0,210,255,0.6);
  }

  .actions{display:flex;gap:14px;justify-content:center;flex-wrap:wrap}
  button{
    cursor:pointer;border:none;border-radius:999px;padding:14px 28px;
    font-size:15px;font-weight:700;letter-spacing:0.03em;
    background:linear-gradient(90deg,var(--accent1),var(--accent2));
    color:#fff;box-shadow:0 10px 30px rgba(106,17,203,0.45);
    transition:transform .15s ease, box-shadow .15s ease;
  }
  button:hover{transform:translateY(-3px) scale(1.04);box-shadow:0 16px 40px rgba(106,17,203,0.6)}
  button:active{transform:translateY(0) scale(0.98)}
  button.alt{background:linear-gradient(90deg,var(--accent3),#00ff99)}
  button.alt2{background:linear-gradient(90deg,var(--gold),var(--accent1))}

  footer{
    position:fixed;left:16px;bottom:14px;font-size:12px;opacity:0.7;z-index:3;
    letter-spacing:0.03em;
  }
  footer b{color:var(--gold)}

  #shake.shaking{animation:shake .5s}
  @keyframes shake{
    10%,90%{transform:translate3d(-1px,0,0)}
    20%,80%{transform:translate3d(2px,0,0)}
    30%,50%,70%{transform:translate3d(-4px,0,0)}
    40%,60%{transform:translate3d(4px,0,0)}
  }

  @media (hover:none){ body{cursor:auto} #cursor{display:none} }
</style>
</head>
<body id="shake">
  <canvas id="stars"></canvas>
  <canvas id="matrix"></canvas>
  <canvas id="confetti"></canvas>
  <canvas id="cursor" width="40" height="40"></canvas>

  <div class="stage">
    <h1 id="title"></h1>
    <p class="sub">de um simples <b>console.log</b> para uma <b>experiência multimídia completa</b>.</p>
    <p class="counter">🌍 sim: <span id="count">0</span></p>
    <div class="actions">
      <button id="boom">🎉 Explodir confete</button>
      <button class="alt" id="shakeBtn">📳 Sacudir tela</button>
      <button class="alt2" id="rainbow">🌈 Modo caótico</button>
    </div>
  </div>

  <footer>by somma with much coffe ☕</footer>

<script>
/* ---------- title letter-splitting ---------- */
const text = "HELLO, WORLD!";
const titleEl = document.getElementById('title');
titleEl.innerHTML = Array.from(text).map(function(c){ return '<span>' + (c === ' ' ? '&nbsp;' : c) + '</span>'; }).join('');

/* ---------- counter ---------- */
let count = 0;
const countEl = document.getElementById('count');
function animateCount(target){
  const start = count, dur = 900, t0 = performance.now();
  function step(t){
    const p = Math.min(1,(t-t0)/dur);
    count = Math.floor(start + (target-start)*(1-Math.pow(1-p,3)));
    countEl.textContent = count.toLocaleString('pt-BR');
    if(p<1) requestAnimationFrame(step);
  }
  requestAnimationFrame(step);
}
animateCount(10000);

/* ---------- confetti ---------- */
const confetti = document.getElementById('confetti');
const cctx = confetti.getContext('2d');
function resizeC(){confetti.width=innerWidth;confetti.height=innerHeight}
addEventListener('resize',resizeC);resizeC();
const parts=[];
function spawnConfetti(x,y,n=180){
  for(let i=0;i<n;i++){
    parts.push({
      x,y,
      vx:(Math.random()-0.5)*14,
      vy:(Math.random()-1.8)*14,
      sz:Math.random()*8+3,
      rot:Math.random()*Math.PI,
      vr:(Math.random()-0.5)*0.3,
      c:'hsl('+Math.floor(Math.random()*360)+' 90% 60%)',
      life:Math.random()*90+50
    });
  }
}
function tickConfetti(){
  cctx.clearRect(0,0,confetti.width,confetti.height);
  for(let i=parts.length-1;i>=0;i--){
    const p=parts[i];
    p.x+=p.vx;p.y+=p.vy;p.vy+=0.35;p.rot+=p.vr;p.life--;
    cctx.save();
    cctx.translate(p.x,p.y);cctx.rotate(p.rot);
    cctx.fillStyle=p.c;
    cctx.fillRect(-p.sz/2,-p.sz/2,p.sz,p.sz);
    cctx.restore();
    if(p.life<=0||p.y>confetti.height+50) parts.splice(i,1);
  }
  requestAnimationFrame(tickConfetti);
}
tickConfetti();
document.getElementById('boom').addEventListener('click',()=>{
  spawnConfetti(innerWidth/2,innerHeight/2,220);
  animateCount(count+Math.floor(Math.random()*5000)+1000);
});

/* ---------- matrix rain ---------- */
const m = document.getElementById('matrix');
const mx = m.getContext('2d');
function mr(){m.width=innerWidth;m.height=innerHeight}
addEventListener('resize',mr);mr();
let cols = Math.floor(m.width/14);
let drops = Array(cols).fill(1);
function matrix(){
  mx.fillStyle='rgba(0,0,0,0.08)';
  mx.fillRect(0,0,m.width,m.height);
  mx.fillStyle='#00ff99';
  mx.font='14px monospace';
  for(let i=0;i<drops.length;i++){
    const ch = String.fromCharCode(0x30A0+Math.random()*96);
    mx.fillText(ch,i*14,drops[i]*14);
    if(drops[i]*14>m.height && Math.random()>0.975) drops[i]=0;
    drops[i]++;
  }
  requestAnimationFrame(matrix);
}
matrix();

/* ---------- starfield ---------- */
const stars = document.getElementById('stars');
const sx = stars.getContext('2d');
function sr(){stars.width=innerWidth;stars.height=innerHeight}
addEventListener('resize',sr);sr();
const starList = Array.from({length:150},()=>({
  x:Math.random()*stars.width,
  y:Math.random()*stars.height,
  r:Math.random()*1.6+0.3,
  tw:Math.random()*0.02+0.005,
  ph:Math.random()*Math.PI*2
}));
function drawStars(t){
  sx.clearRect(0,0,stars.width,stars.height);
  for(const s of starList){
    const a = 0.5+0.5*Math.sin(t*s.tw+s.ph);
    sx.beginPath();
    sx.fillStyle='rgba(255,255,255,'+a+')';
    sx.arc(s.x,s.y,s.r,0,Math.PI*2);
    sx.fill();
  }
  requestAnimationFrame(drawStars);
}
requestAnimationFrame(drawStars);

/* ---------- custom cursor with trail ---------- */
const cursor = document.getElementById('cursor');
const cx2 = cursor.getContext('2d');
let mouseX=innerWidth/2, mouseY=innerHeight/2;
const trail=[];
addEventListener('mousemove',e=>{
  mouseX=e.clientX;mouseY=e.clientY;
  trail.push({x:mouseX,y:mouseY,life:20});
  cursor.style.left='0px';cursor.style.top='0px';
});
cursor.width=innerWidth;cursor.height=innerHeight;
addEventListener('resize',()=>{cursor.width=innerWidth;cursor.height=innerHeight});
function drawCursor(){
  cx2.clearRect(0,0,cursor.width,cursor.height);
  for(let i=trail.length-1;i>=0;i--){
    const p=trail[i];p.life--;
    const a=p.life/20;
    cx2.beginPath();
    cx2.fillStyle='hsla('+((i*20)%360)+',90%,65%,'+(a*0.6)+')';
    cx2.arc(p.x,p.y,6*a,0,Math.PI*2);
    cx2.fill();
    if(p.life<=0) trail.splice(i,1);
  }
  cx2.beginPath();
  cx2.fillStyle='#fff';
  cx2.arc(mouseX,mouseY,4,0,Math.PI*2);
  cx2.fill();
  requestAnimationFrame(drawCursor);
}
drawCursor();

/* ---------- shake button ---------- */
const body = document.getElementById('shake');
document.getElementById('shakeBtn').addEventListener('click',()=>{
  body.classList.remove('shaking');
  void body.offsetWidth;
  body.classList.add('shaking');
  spawnConfetti(innerWidth/2,innerHeight*0.3,80);
});

/* ---------- chaos mode: randomize accent colors ---------- */
const root = document.documentElement.style;
document.getElementById('rainbow').addEventListener('click',()=>{
  const rand = function(){ return 'hsl('+Math.floor(Math.random()*360)+' 90% 60%)'; };
  root.setProperty('--accent1',rand());
  root.setProperty('--accent2',rand());
  root.setProperty('--accent3',rand());
  root.setProperty('--gold',rand());
  spawnConfetti(innerWidth/2,innerHeight/2,150);
});

/* ---------- click anywhere = tiny burst ---------- */
addEventListener('click',e=>{
  if(e.target.tagName!=='BUTTON') spawnConfetti(e.clientX,e.clientY,30);
});
</script>
</body>
</html>

`

    c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	})

	log.Fatal(app.Listen(":3000"))
}
