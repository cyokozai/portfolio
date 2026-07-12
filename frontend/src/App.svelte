<script>
  import { fly } from 'svelte/transition';

  function onVisible(node, callback) {
    const observer = new IntersectionObserver(
      ([entry]) => { if (entry.isIntersecting) callback(); },
      { threshold: 0.4 }
    );
    observer.observe(node);
    return { destroy: () => observer.disconnect() };
  }

  let mottosVisible = false;
</script>

<!-- NAV -->
<nav>
  <a class="nav-brand" href="/">cyokozai</a>
</nav>
<div class="rule rule--cyan"></div>

<!-- HERO -->
<section id="hero">
  <div class="hero-content">
    <p class="hero-label">Platform Engineer · Graduate Student</p>
    <h1>Yusuke Inoue</h1>
    <p class="hero-sub">a.k.a. cyokozai | 猪口才</p>
    <p class="hero-tagline"></p>
    <ul class="hero-affil">
      <li>千葉工業大学大学院 情報科学研究科 情報科学専攻</li>
      <li>CloudNative Days 実行委員</li>
      <li>アメリス株式会社 インターン</li>
    </ul>
    <div class="hero-mottos">
      <span class="motto">Infrastructure as a Smile</span>
      <span class="motto">今ハタダ、深ク静カニ潜航セヨ。</span>
    </div>
  </div>
  <div class="hero-photo" in:fly={{ y: 16, duration: 900, delay: 300 }}>
    <img src="https://github.com/cyokozai.png" alt="cyokozai" />
  </div>
</section>

<!-- MOTTOS -->
<section id="mottos">
  <div class="cyan-block" use:onVisible={() => mottosVisible = true}>
    {#if mottosVisible}
      <p class="block-motto" in:fly={{ y: 20, duration: 600 }}>Infrastructure as a Smile</p>
      <p class="block-motto" in:fly={{ y: 20, duration: 600, delay: 150 }}>今ハタダ、深ク静カニ潜航セヨ。</p>
    {/if}
  </div>
</section>

<!-- FOOTER -->
<footer>
  <p>© 2026 Yusuke Inoue</p>
</footer>


<style>
  :global(#app) {
    width: 760px;
    max-width: 100%;
    margin: 0 auto;
    border: 1px solid var(--border);
    min-height: 100dvh;
    display: flex;
    flex-direction: column;
  }

  /* ── Accent rules ── */
  .rule {
    flex-shrink: 0;
    width: 100%;
  }

  .rule--cyan {
    height: 3px;
    background: var(--accent);
  }

  .rule--yellow {
    height: 4px;
    background: var(--accent-sub);
  }

  /* ── Nav ── */
  nav {
    display: flex;
    align-items: center;
    padding: 22px 40px;
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    background: var(--bg);
    z-index: 10;
  }

  .nav-brand {
    font-size: 15px;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--text-h);
  }

  /* ── Hero ── */
  #hero {
    display: grid;
    grid-template-columns: 1fr 200px;
    gap: 32px;
    align-items: start;
    padding: 72px 40px 80px;
  }

  .hero-label {
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 1.2px;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 20px;
  }

  #hero h1 {
    font-size: 52px;
    font-weight: 700;
    letter-spacing: -1.5px;
    color: var(--text-h);
    line-height: 1.05;
    margin: 0 0 14px;
  }

  .hero-sub {
    font-size: 16px;
    font-weight: 400;
    color: var(--text);
    margin-bottom: 28px;
    letter-spacing: 0.2px;
  }

  .hero-tagline {
    font-size: 15px;
    line-height: 1.8;
    color: var(--text);
    margin-bottom: 8px;
  }

  .hero-affil {
    list-style: none;
    padding: 0;
    margin: 0 0 36px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .hero-affil li {
    font-size: 13px;
    color: var(--text);
    letter-spacing: 0.05em;
    opacity: 0.75;
  }

  .hero-affil li::before {
    content: '·  ';
    color: var(--accent);
    font-weight: 700;
  }

  .hero-mottos {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
    
  }

  .motto {
    font-family: var(--mono);
    font-size: 13px;
    color: var(--accent);
    background: var(--accent-bg);
    padding: 5px 14px;
    border-radius: 4px;
    border: 1px solid var(--accent-border);
    letter-spacing: 0.2px;
  }

  .hero-photo {
    width: 200px;
    aspect-ratio: 1/1;
    border-radius: 8px;
    overflow: hidden;
    margin-top: 8px;
  }

  .hero-photo img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center top;
    display: block;
  }

  .img-label {
    font-size: 11px;
    font-family: var(--mono);
    color: var(--border);
    letter-spacing: 0.5px;
    text-transform: uppercase;
  }

  /* ── Mottos block ── */
  #mottos {
    padding: 80px 40px;
    display: flex;
    justify-content: center;
  }

  .cyan-block {
    width: 200px;
    height: 200px;
    background: var(--color-cyan);
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    padding: 24px;
    text-align: center;
  }

  .block-motto {
    font-size: 13px;
    font-family: var(--mono);
    color: var(--color-white);
    letter-spacing: 0.05em;
    line-height: 1.6;
  }

  /* ── Footer ── */
  footer {
    padding: 24px 40px;
    border-top: 4px solid var(--accent-sub);
  }

  footer p {
    font-size: 13px;
    color: var(--text);
  }

  /* ── Responsive ── */
  @media (max-width: 600px) {
    nav {
      padding: 18px 20px;
    }

    #hero {
      grid-template-columns: 1fr;
      padding: 48px 24px 56px;
    }

    #hero h1 {
      font-size: 38px;
    }

    .hero-photo {
      width: 160px;
      aspect-ratio: 1/1;
    }

    footer {
      padding: 20px 24px;
    }
  }
</style>
