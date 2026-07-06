<script>
  const navLinks = [
    { label: 'Projects', href: '#projects' },
    { label: 'Writing', href: '#writing' },
    { label: 'Talks', href: '#talks' },
    { label: 'Community', href: '#community' },
    { label: 'Contact', href: '#contact' },
  ]

  const community = [
    {
      role: '実行委員 + 会場DJ',
      org: 'CloudNative Days',
      desc: 'Cloud Native の国内最大コミュニティイベントの運営スタッフ。当日は会場でターンテーブルも回す。',
      link: 'https://cloudnativedays.jp',
    },
    {
      role: 'NOC Server チーム',
      org: 'NeGI NOC',
      desc: 'カンファレンス（IWSC Shizuoka 等）のネットワークインフラを現地で構築・運用。Prometheus / Grafana でリアルタイム監視。',
      link: 'https://github.com/cyokozai',
    },
    {
      role: 'Nekko Cloud Team',
      org: 'ネットワークコンテンツ研究会',
      desc: 'サークルのプライベートクラウド "Nekko Cloud" を運用。仮想マシン操作を自動化する ChatOps ツールをゼロから設計・実装。',
      link: 'https://github.com/cyokozai',
    },
  ]

  const socials = [
    { label: 'GitHub', href: 'https://github.com/cyokozai' },
    { label: 'SpeakerDeck', href: 'https://speakerdeck.com/yusuke427' },
    { label: 'Zenn', href: 'https://zenn.dev/cyokozai' },
    { label: 'Qiita', href: 'https://qiita.com/cyokozai' },
  ]

  const skills = [
    'Go', 'Python', 'Julia', 'Kubernetes', 'Terraform',
    'Docker', 'GCP', 'eBPF', 'CI/CD', 'Linux', 'Svelte',
  ]

  const projects = [
    {
      title: 'Terraform × cloud-init VMセットアップ自動化',
      desc: 'Infrastructure as Code でのVM自動構成。登壇スライド 2.6k 閲覧。',
      tags: ['Terraform', 'cloud-init', 'IaC'],
      link: 'https://speakerdeck.com/yusuke427',
    },
    {
      title: 'Platform Engineering 視点のK8s LB比較',
      desc: 'Standalone NEG / Ingress / Gateway API を所掌分離・認知負荷・ヘルスチェックの軸で比較検証。',
      tags: ['Kubernetes', 'GKE', 'Platform Engineering'],
      link: 'https://sreake.com/blog/kubernetes-load-balancing-comparison-on-platform-engineering-view-point/',
    },
    {
      title: 'Tetragon vs PodSecurity/seccomp 比較',
      desc: 'eBPF ベースのランタイムセキュリティと K8s 標準セキュリティ機能を実践的に比較。スリーシェイク インターン成果。',
      tags: ['eBPF', 'Kubernetes', 'Security'],
      link: 'https://sreake.com/blog/container-security-comparison/',
    },
  ]

  const writings = [
    {
      title: 'Cloud Service Mesh 入門',
      desc: 'Google Cloud マネージドサービスメッシュの構築と Istio との比較。',
      source: 'Sreake Blog',
      link: 'https://sreake.com/blog/cloud-service-mesh-getting-started/',
    },
    {
      title: 'メタヒューリスティクスを Rust/WASM で実装する',
      desc: '最適化アルゴリズムをブラウザで動かす実験。3-shake Advent Calendar 2025。',
      source: 'Qiita',
      link: 'https://qiita.com/cyokozai/items/93b41ca1d5407e6537ff',
    },
    {
      title: '若手・学生がテックカンファレンスに参加すべき5つの理由',
      desc: 'CloudNative Days 実行委員として書いた学生向けガイド。',
      source: 'CloudNative Days Blog',
      link: 'https://kaigi.cloudnativedays.jp/blog/conference-tips-for-juniors/',
    },
  ]

  const talks = [
    { title: 'シェーダーもGoで書く！Ebitengineで作るリアルタイムGPUシェーダーエディタ', event: '---' },
    { title: '年末なのでオンプレKubernetesクラスタをアップグレードしよう', event: '---' },
    { title: 'KubeVirtのLive Migration実践', event: '---' },
    { title: 'Terraform × cloud-initでVMセットアップ自動化', event: '---' },
    { title: '10Gbルータを作るための下準備', event: '---' },
  ]
</script>

<!-- NAV -->
<nav>
  <a class="nav-brand" href="/">cyokozai</a>
  <ul class="nav-links">
    {#each navLinks as { label, href }}
      <li><a {href}>{label}</a></li>
    {/each}
  </ul>
</nav>

<!-- HERO -->
<section id="hero">
  <p class="hero-label">Platform Engineer · Graduate Student · TA</p>
  <h1>cyokozai | 猪口才</h1>
  <h2>a.k.a. Yusuke Inoue</h2>
  <p class="hero-tagline">
    Cloud Native な基盤を設計・構築し、プロダクトとして届ける。<br>
    千葉工業大学大学院 &nbsp;·&nbsp; CloudNative Days 実行委員 &nbsp;·&nbsp; スリーシェイク インターン
  </p>
  <div class="hero-mottos">
    <span class="motto-code">Infrastructure as a Smile</span>
    <span class="motto-sep">·</span>
    <span class="motto-ja">今ハタダ、深ク静カニ潜航セヨ。</span>
  </div>
  <div class="hero-socials">
    {#each socials as { label, href }}
      <a class="pill" {href} target="_blank" rel="noopener noreferrer">{label} ↗</a>
    {/each}
  </div>
</section>

<div class="divider"></div>

<!-- SKILLS -->
<section id="skills">
  <h2>Skills</h2>
  <div class="chips">
    {#each skills as skill}
      <span class="chip">{skill}</span>
    {/each}
  </div>
</section>

<div class="divider"></div>

<!-- PROJECTS -->
<section id="projects">
  <h2>Projects</h2>
  <p class="section-desc">思考プロセスと成果を残したもの</p>
  <div class="card-grid">
    {#each projects as { title, desc, tags, link }}
      <a class="card" href={link} target="_blank" rel="noopener noreferrer">
        <h3>{title}</h3>
        <p>{desc}</p>
        <div class="tags">
          {#each tags as tag}
            <span class="tag">{tag}</span>
          {/each}
        </div>
      </a>
    {/each}
  </div>
</section>

<div class="divider"></div>

<!-- WRITING -->
<section id="writing">
  <h2>Writing</h2>
  <p class="section-desc">技術記事・アドベントカレンダー・コミュニティブログ</p>
  <div class="writing-list">
    {#each writings as { title, desc, source, link }}
      <a class="writing-item" href={link} target="_blank" rel="noopener noreferrer">
        <span class="writing-source">{source}</span>
        <h3>{title}</h3>
        <p>{desc}</p>
      </a>
    {/each}
  </div>
</section>

<div class="divider"></div>

<!-- TALKS -->
<section id="talks">
  <h2>Talks</h2>
  <p class="section-desc">登壇資料（SpeakerDeck に全 14 本公開）</p>
  <ul class="talks-list">
    {#each talks as { title, event }}
      <li>
        <span class="talks-dot">·</span>
        <span>{title}</span>
      </li>
    {/each}
    <li class="talks-more">
      <a href="https://speakerdeck.com/yusuke427" target="_blank" rel="noopener noreferrer">
        すべての登壇資料を見る →
      </a>
    </li>
  </ul>
</section>

<div class="divider"></div>

<!-- COMMUNITY -->
<section id="community">
  <h2>Community</h2>
  <p class="section-desc">技術コミュニティでの活動・役割</p>
  <div class="community-grid">
    {#each community as { role, org, desc, link }}
      <a class="community-card" href={link} target="_blank" rel="noopener noreferrer">
        <div class="community-header">
          <span class="community-org">{org}</span>
          <span class="community-role">{role}</span>
        </div>
        <p>{desc}</p>
      </a>
    {/each}
  </div>
</section>

<div class="divider"></div>

<!-- CONTACT -->
<section id="contact">
  <h2>Contact</h2>
  <p class="contact-desc">
    コラボレーション・案件相談・登壇依頼はお気軽にどうぞ。<br>
    GitHub の DM か Issues でも歓迎です。
  </p>
  <a class="contact-cta" href="https://github.com/cyokozai" target="_blank" rel="noopener noreferrer">
    GitHub で連絡する →
  </a>
</section>

<!-- FOOTER -->
<footer>
  <p>© 2026 Yusuke Inoue</p>
</footer>

<style>
  /* ── Layout ── */
  :global(#app) {
    width: 760px;
    max-width: 100%;
    margin: 0 auto;
    border-inline: 1px solid var(--border);
    min-height: 100svh;
    display: flex;
    flex-direction: column;
  }

  /* ── Nav ── */
  nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 40px;
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    background: var(--bg);
    z-index: 10;
  }

  .nav-brand {
    font-weight: 600;
    color: var(--text-h);
    font-size: 15px;
    letter-spacing: 0.5px;
  }

  .nav-links {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 28px;
  }

  .nav-links a {
    font-size: 14px;
    color: var(--text);
    transition: color 0.2s;
  }

  .nav-links a:hover {
    color: var(--text-h);
  }

  /* ── Sections ── */
  section {
    padding: 64px 40px;
  }

  h2 {
    font-size: 13px;
    font-weight: 600;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 16px;
  }

  .section-desc {
    font-size: 14px;
    color: var(--text);
    margin-bottom: 32px;
  }

  .divider {
    height: 1px;
    background: var(--border);
    margin: 0 40px;
  }

  /* ── Hero ── */
  #hero {
    padding-top: 80px;
    padding-bottom: 80px;
  }

  .hero-label {
    font-size: 13px;
    font-weight: 500;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 16px;
  }

  #hero h1 {
    font-size: 48px;
    font-weight: 700;
    letter-spacing: -1.5px;
    color: var(--text-h);
    line-height: 1.1;
    margin: 0 0 20px;
  }

  .hero-tagline {
    font-size: 17px;
    line-height: 1.7;
    color: var(--text);
    margin-bottom: 36px;
    max-width: 560px;
  }

  .hero-mottos {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
    margin-bottom: 32px;
  }

  .motto-code {
    font-family: var(--mono);
    font-size: 14px;
    color: var(--accent);
    background: var(--accent-bg);
    padding: 4px 12px;
    border-radius: 6px;
    border: 1px solid var(--accent-border);
  }

  .motto-sep {
    color: var(--border);
    font-size: 18px;
  }

  .motto-ja {
    font-size: 14px;
    color: var(--text);
    letter-spacing: 0.05em;
  }

  .hero-socials {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .pill {
    display: inline-flex;
    align-items: center;
    padding: 7px 16px;
    border-radius: 100px;
    border: 1px solid var(--border);
    font-size: 14px;
    color: var(--text-h);
    background: var(--bg);
    transition: border-color 0.2s, background 0.2s;
  }

  .pill:hover {
    border-color: var(--accent-border);
    background: var(--accent-bg);
  }

  /* ── Skills ── */
  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .chip {
    padding: 5px 14px;
    border-radius: 100px;
    background: var(--code-bg);
    font-size: 13px;
    font-family: var(--mono);
    color: var(--text-h);
  }

  /* ── Projects ── */
  .card-grid {
    display: grid;
    gap: 16px;
  }

  .card {
    display: block;
    padding: 24px;
    border: 1px solid var(--border);
    border-radius: 12px;
    transition: border-color 0.2s, box-shadow 0.2s;
  }

  .card:hover {
    border-color: var(--accent-border);
    box-shadow: 0 4px 20px rgba(170, 59, 255, 0.06);
  }

  .card h3 {
    font-size: 15px;
    font-weight: 600;
    margin-bottom: 8px;
    color: var(--text-h);
  }

  .card p {
    font-size: 14px;
    line-height: 1.6;
    color: var(--text);
    margin-bottom: 16px;
  }

  .tags {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
  }

  .tag {
    font-size: 12px;
    padding: 3px 10px;
    border-radius: 100px;
    background: var(--accent-bg);
    color: var(--accent);
    font-weight: 500;
  }

  /* ── Writing ── */
  .writing-list {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .writing-item {
    display: block;
    padding: 20px 0;
    border-bottom: 1px solid var(--border);
    transition: opacity 0.2s;
  }

  .writing-item:first-child {
    border-top: 1px solid var(--border);
  }

  .writing-item:hover {
    opacity: 0.75;
  }

  .writing-source {
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.8px;
    text-transform: uppercase;
    color: var(--accent);
  }

  .writing-item h3 {
    font-size: 15px;
    font-weight: 600;
    margin: 6px 0 6px;
    color: var(--text-h);
  }

  .writing-item p {
    font-size: 14px;
    color: var(--text);
  }

  /* ── Talks ── */
  .talks-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .talks-list li {
    display: flex;
    gap: 12px;
    font-size: 15px;
    color: var(--text-h);
    align-items: baseline;
  }

  .talks-dot {
    color: var(--accent);
    font-size: 20px;
    line-height: 1;
    flex-shrink: 0;
  }

  .talks-more a {
    font-size: 14px;
    color: var(--accent);
    font-weight: 500;
  }

  .talks-more a:hover {
    text-decoration: underline;
  }

  /* ── Community ── */
  .community-grid {
    display: grid;
    gap: 16px;
  }

  .community-card {
    display: block;
    padding: 20px 24px;
    border: 1px solid var(--border);
    border-radius: 12px;
    transition: border-color 0.2s, background 0.2s;
  }

  .community-card:hover {
    border-color: var(--accent-border);
    background: var(--accent-bg);
  }

  .community-header {
    display: flex;
    align-items: baseline;
    gap: 12px;
    margin-bottom: 10px;
    flex-wrap: wrap;
  }

  .community-org {
    font-size: 15px;
    font-weight: 700;
    color: var(--text-h);
  }

  .community-role {
    font-size: 12px;
    font-weight: 600;
    padding: 3px 10px;
    border-radius: 100px;
    background: var(--accent-bg);
    color: var(--accent);
    border: 1px solid var(--accent-border);
  }

  .community-card p {
    font-size: 14px;
    line-height: 1.65;
    color: var(--text);
  }

  /* ── Contact ── */
  .contact-desc {
    font-size: 16px;
    line-height: 1.7;
    color: var(--text);
    margin-bottom: 28px;
  }

  .contact-cta {
    display: inline-flex;
    align-items: center;
    padding: 12px 24px;
    background: var(--accent);
    color: #fff;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 500;
    transition: opacity 0.2s;
  }

  .contact-cta:hover {
    opacity: 0.85;
  }

  /* ── Footer ── */
  footer {
    padding: 32px 40px;
    border-top: 1px solid var(--border);
    margin-top: auto;
  }

  footer p {
    font-size: 13px;
    color: var(--text);
  }

  /* ── Responsive ── */
  @media (max-width: 600px) {
    nav {
      padding: 16px 20px;
    }

    .nav-links {
      gap: 16px;
    }

    section {
      padding: 48px 20px;
    }

    .divider {
      margin: 0 20px;
    }

    #hero h1 {
      font-size: 36px;
    }
  }
</style>
