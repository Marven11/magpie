// Who magpie's releases came from: the people a commit credits, read from
// what is already public in yetone/magpie's history. Nothing here is made up
// and no email is ever given out:
//
//   - a subject's parenthesised credit: "(hoobnn on Discord)",
//     "(qikamatata on X: …)", "(#2356, detail47)", "(willz, Discord)",
//     "(TJHHHH, luci on Discord)", a lone "(mrwl666)";
//   - the author of an issue or PR such a group names, "(#4532)" (looked up
//     by credits-cli.mjs; a lookup that fails names nobody);
//   - a commit's GitHub author and Co-authored-by names other than the
//     maintainer's and bots'.
//
// A lone word in brackets is often code ("(nogui)", "(P3)", "(Codex)"), so it
// counts only when it is shaped like a handle with digits ("mrwl666") or the
// same name is credited plainly somewhere else ("axelburks on Discord").
//
// Used by credits-cli.mjs for the release notes' thanks line. It imports
// nothing, so the site's worker can use it too.

export const MAINTAINERS = ["yetone"];
// agents and bots that author or co-author commits
const BOT = /\[bot\]$|^(claude( code| opus.*| sonnet.*| haiku.*)?|amp|codex|copilot|github copilot|cursor( agent)?|devin( ai)?|gemini|opencode|droid|factory-droid|dependabot|github-actions|renovate|jev)$/i;
// agents and machine accounts that co-author commits: "Instinct Agent
// <agent@instinct.com>", "glm-5.3-flash", "magpie-fix <fix@localhost>"
const AGENT_MAIL = /@(anthropic\.com|cursor\.com|openai\.com|instinct\.com|multica\.ai|factory\.ai|ampcode\.com|google\.com|localhost)$|^(noreply|agent|bot)@/i;
const MODEL = /^(claude|gpt|glm|qwen|kimi|deepseek|gemini|grok|minimax|doubao|o\d)\b.*\d|\bagent$|^(admin|root|user|ubuntu|runner|test)$/i;
const PLATFORMS = { discord: "Discord", x: "X", twitter: "X", github: "GitHub", v2ex: "V2EX", wechat: "WeChat", 微信: "WeChat", telegram: "Telegram", reddit: "Reddit", "linux.do": "LINUX DO", 即刻: "Jike", 小红书: "Xiaohongshu" };
// words that sit in brackets beside an issue number and aren't people
const STOP = new Set(("discord x github zh en ja de zh-tw zh-cn high low medium ultra xhigh max min mini pro plus free " +
  "model models folder download review fixes fixed see also and or not windows linux macos mac wsl ios android ci " +
  "main red green beta alpha stable new old test tests api cli gui tui web site docs owner user users the pi spelled " +
  "plugin plugins provider providers agent agents gateway routing usage settings panel tray").split(" "));
const LOGIN = /^[A-Za-z0-9](?:[A-Za-z0-9]|-(?=[A-Za-z0-9])){2,38}$/;
// a Discord-style handle ending in its number: tjhhhh_84378
const NUMBERED = /^_?[A-Za-z][A-Za-z0-9]*(?:_[A-Za-z0-9]+)*_\d{3,}$/;
// a lone bracketed word taken without anything else naming it: mrwl666
const DIGITED = /^[A-Za-z]{3,}\d{2,}$/;

const lower = (s) => s.toLowerCase();

// a word next to an issue number or another credit that may be a person
function handle(w) {
  if (STOP.has(lower(w))) return false;
  if (/^\d+$/.test(w) || /^[0-9a-f]{7,40}$/.test(w) || /^v?\d+(\.\d+)+/.test(w)) return false;
  return LOGIN.test(w) || NUMBERED.test(w);
}

// someone named in front of a platform, "Benjie Zhou on Discord"
function personName(s) {
  s = s.trim().replace(/^@/, "");
  if (!s || s.length > 40 || /[#()\[\]{}<>"“”]/.test(s) || s.split(/\s+/).length > 3) return "";
  if (/^(the|a|an)\s/i.test(s) || STOP.has(lower(s))) return "";
  // who someone is, not their name: "Docker user", "𝕏 user", "𝕏"
  if (/(^|\s)(user|users|someone|reporter|people)$/i.test(s) || !/[\p{L}\p{N}]/u.test(s.replace(/𝕏/gu, ""))) return "";
  return s;
}

// The credits in one subject: names (sure, or to be taken only if seen
// elsewhere) and the issues or PRs its brackets name.
export function parse(subject) {
  // older subjects carry their whole body after " — "; credits come before
  const head = String(subject).split(" — ")[0];
  const out = { names: [], refs: [] };
  for (const m of head.matchAll(/\(([^()]*)\)/g)) group(m[1], out);
  return out;
}

function group(text, out) {
  const bare = [];
  let parts = 0;
  for (let piece of text.split(/;\s*/)) {
    // a quote after the name: "NAME on X: 打开页面很卡", "#538: …"
    const colon = piece.search(/:(\s|["“「]|$)/);
    if (colon >= 0) piece = piece.slice(0, colon);
    // a piece names someone only when its own words do: "(#615: …; Neon,
    // Notion and Supabase stay)" doesn't
    let context = false;
    const mine = [];
    for (let p of piece.split(/,\s*/)) {
      p = p.trim();
      if (!p) continue;
      parts++;
      // "spotted by @hezhizhen in #103"
      const by = p.match(/(?:^|\s)by @([A-Za-z0-9-]{1,39})\b/);
      if (by) { out.names.push({ name: by[1], via: "GitHub", sure: true }); context = true; }
      // "#12", "plugins#12", "plugins #12", "magpie-community/plugins#12"
      p = p.replace(/(?:\b(plugins)\s+|([\w.-]+\/)?([\w.-]+))?#(\d+)\b/g, (_, spaced, owner, repo, n) => {
        const r = spaced || (!repo ? "yetone/magpie" : `${owner || ""}${repo}`);
        const known = { "yetone/magpie": "yetone/magpie", magpie: "yetone/magpie", plugins: "magpie-community/plugins", "magpie-community/plugins": "magpie-community/plugins" }[r];
        if (known) out.refs.push({ repo: known, n: +n });
        context = true;
        return " ";
      }).trim();
      if (!p || by) continue;
      const plat = PLATFORMS[lower(p)];
      if (plat) { // "willz, Discord"
        if (mine.length) { out.names.push({ name: mine.pop(), via: plat, sure: true }); context = true; }
        continue;
      }
      const on = p.match(/^(.+?)\s+on\s+(\S+)$/);
      if (on && PLATFORMS[lower(on[2])]) {
        for (const n of on[1].split(/\s+(?:and|&)\s+/)) {
          const name = personName(n);
          if (name) out.names.push({ name, via: PLATFORMS[lower(on[2])], sure: true });
        }
        context = true;
        continue;
      }
      if (/^@?\S+$/.test(p)) mine.push(p.replace(/^@/, ""));
    }
    for (const b of mine) bare.push({ b, context });
  }
  // a lone word, or one beside a credit or a number ("July, LCorleone",
  // "memory, Osal on Discord", "#860 flake", "#2606 fatcutegg"), is a person
  // only if named plainly elsewhere or shaped like a handle with digits:
  // credit() decides
  for (const { b, context } of bare) {
    if (context ? handle(b) : parts === 1 && (LOGIN.test(b) || NUMBERED.test(b))) out.names.push({ name: b, via: null, sure: false });
  }
}

// the GitHub login in a noreply address; the address itself is never kept
function noreply(email) {
  const m = String(email || "").match(/^(?:\d+\+)?([A-Za-z0-9-]+)@users\.noreply\.github\.com$/i);
  return m ? m[1] : "";
}

export function refKey(r) { return `${r.repo}#${r.n}`; }

// excluded: maintainer, bots, the people who asked to be left off
function left(name, exclude) {
  const k = lower(name);
  return !k || MAINTAINERS.includes(k) || BOT.test(name) || exclude.has(k) || /^(the |a )?(owner|user|users|maintainer)$/.test(k);
}

// a co-author or commit author who is a person, not an agent or machine
function human(name, email) {
  if (!name || BOT.test(name) || MODEL.test(name) || !/^[\p{L}\p{N}]/u.test(name)) return false;
  return !email || !AGENT_MAIL.test(email) || /@users\.noreply\.github\.com$/i.test(email);
}
const coName = (c) => (typeof c === "string" ? { name: c, email: "" } : c);

// Each commit's people. commits: [{ subject, author, email, parents,
// coauthors: [{ name, email }] }]
// issueAuthors: Map or object of "repo#n" → login (a bot's ends in "[bot]").
// Returns an array beside commits, each [{ key, name, via, github }].
export function credit(commits, { issueAuthors = {}, exclude = [] } = {}) {
  const authorOf = issueAuthors instanceof Map ? (k) => issueAuthors.get(k) : (k) => issueAuthors[k];
  const ex = new Set([...exclude].map((s) => lower(String(s).trim())).filter(Boolean));
  const parsed = commits.map((c) => parse(c.subject));
  // names credited plainly anywhere, so a lone bracketed one can be trusted
  const known = new Set();
  commits.forEach((c, i) => {
    for (const n of parsed[i].names) if (n.sure) known.add(lower(n.name));
    for (const r of parsed[i].refs) { const a = authorOf(refKey(r)); if (a) known.add(lower(a)); }
    if (human(c.author, c.email)) known.add(lower(c.author));
    const gh = noreply(c.email);
    if (gh) known.add(lower(gh));
    for (const co of (c.coauthors || []).map(coName)) if (human(co.name, co.email)) known.add(lower(co.name));
  });
  return commits.map((c, i) => {
    const people = new Map();
    const add = (name, via, github) => {
      name = String(name || "").trim();
      if (left(name, ex) || (github && left(github, ex))) return;
      const key = lower(github || name);
      const had = people.get(key);
      if (had) { if (via && !had.via.includes(via)) had.via.push(via); if (github) had.github = github; return; }
      people.set(key, { key, name, via: via ? [via] : [], github: github || "" });
    };
    for (const n of parsed[i].names) {
      if (n.sure || known.has(lower(n.name)) || (DIGITED.test(n.name) && !/^[0-9a-f]+$/i.test(n.name)) || NUMBERED.test(n.name)) add(n.name, n.via, "");
    }
    for (const r of parsed[i].refs) {
      const a = authorOf(refKey(r));
      if (a && !BOT.test(a)) add(a, "GitHub", a);
    }
    if (c.parents == null || c.parents < 2) {
      if (human(c.author, c.email)) add(c.author, "GitHub", noreply(c.email));
      for (const co of (c.coauthors || []).map(coName)) if (human(co.name, co.email)) add(co.name, "GitHub", noreply(co.email));
    }
    return [...people.values()];
  });
}

// every person across commits, merged by key
export function everyone(perCommit) {
  const all = new Map();
  for (const list of perCommit) for (const p of list) {
    const had = all.get(p.key);
    if (!had) { all.set(p.key, { ...p, via: [...p.via] }); continue; }
    for (const v of p.via) if (!had.via.includes(v)) had.via.push(v);
    if (p.github && !had.github) had.github = p.github;
  }
  return [...all.values()];
}

// markdown that shows a name as written: tjhhhh_84378, not tjhhhh84378 in italics
const md = (s) => s.replace(/([\\`*_{}\[\]<>()#+!|~])/g, "\\$1");

export const WALL = { en: "https://usemagpie.ai/docs/contributors", zh: "https://usemagpie.ai/docs/zh/contributors" };

// the release notes' last line, or "" when nobody is credited
export function thanks(people, lang) {
  const names = people.map((p) => md(p.name));
  if (!names.length) return "";
  if (lang === "zh") {
    const list = names.length > 1 ? names.slice(0, -1).join("、") + " 和 " + names.at(-1) : names[0];
    return `感谢 ${list}：这次更新里有${names.length > 1 ? "你们" : "你"}的反馈与代码。所有贡献者见[贡献墙](${WALL.zh})。`;
  }
  const list = names.length > 1 ? names.slice(0, -1).join(", ") + " and " + names.at(-1) : names[0];
  return `Thanks to ${list}, whose reports and code went into this release. Everyone who helped: [contributors](${WALL.en}).`;
}
