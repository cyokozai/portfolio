<script>
  import { fly } from 'svelte/transition';

  // 本番では Go の単一バイナリが同一オリジンで返す。
  // 開発時は Vite の proxy が backend コンテナへ中継する。
  const profile = fetch('/api/profile').then((res) => {
    if (!res.ok) throw new Error(`プロフィールの取得に失敗しました (${res.status})`);
    return res.json();
  });
</script>

<!-- NAV -->
<nav>
  <a class="nav-brand" href="/">cyokozai</a>
</nav>
<div class="rule rule--cyan"></div>

<!-- HERO -->
{#await profile}
  <section id="hero" aria-busy="true">
    <div class="hero-content">
      <p class="hero-status">読み込み中…</p>
    </div>
  </section>
{:then p}
  <section id="hero">
    <div class="hero-content">
      <p class="hero-label">{p.label}</p>
      <h1>{p.name}</h1>
      <p class="hero-sub">a.k.a. {p.alias}</p>
      <ul class="hero-affil">
        {#each p.affiliations as affiliation}
          <li>{affiliation}</li>
        {/each}
      </ul>
    </div>
    <div class="hero-photo" in:fly={{ y: 16, duration: 900, delay: 300 }}>
      <img src={p.avatarUrl} alt={p.name} />
    </div>
  </section>
{:catch error}
  <section id="hero">
    <div class="hero-content">
      <p class="hero-status hero-status--error" role="alert">{error.message}</p>
    </div>
  </section>
{/await}

<!-- FOOTER -->
<footer>
  <p>© 2026 cyokozai</p>
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

  .hero-status {
    font-size: 14px;
    color: var(--text);
    opacity: 0.6;
  }

  .hero-status--error {
    color: var(--accent-dark);
    opacity: 1;
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
