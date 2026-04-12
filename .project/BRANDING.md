# Rampart — Branding & Identity

> **Version:** 1.0.0  
> **Date:** 2026-04-11

---

## 1. Name & Etymology

**Rampart** — A defensive wall of a castle or walled city, having a broad top with a walkway.

- **Etymology:** Middle French *rempart*, from *remparer* "to fortify"
- **Metaphor:** Rampart stands between your network and threats — a programmable, intelligent wall that remembers everything and never sleeps
- **Sur duvarı** (Türkçe karşılığı)

---

## 2. Writing Style

| Context | Usage |
|---------|-------|
| Title case | Rampart |
| Code / CLI | `rampart` |
| Go module | `github.com/rampartfw/rampart` |
| Config file | `rampart.yaml` |
| Env prefix | `RAMPART_` |
| Docker image | `rampartfw/rampart` |
| Systemd service | `rampart.service` |
| Data directory | `/var/lib/rampart/` |
| Config directory | `/etc/rampart/` |
| Hashtag | #rampart #rampartfw |
| npm scope (WebUI) | `@rampart/ui` |

**Never:** Ram-Part, Ram Part, ram_part, RAMPART (in prose)

---

## 3. Taglines

| Context | Tagline |
|---------|---------|
| **Primary** | Policy-as-Code Firewall. One Binary. Every Backend. |
| **Technical** | The wall that remembers everything. |
| **Developer** | Stop managing iptables rules by hand. |
| **Comparison** | Terraform plan, but for firewalls. |
| **Short** | Programmable firewall. Zero dependencies. |
| **Turkish** | Tek binary. Tüm backend'ler. Tüm kurallar. |

---

## 4. Color Palette

### Primary Colors

| Name | Hex | Usage |
|------|-----|-------|
| **Rampart Slate** | `#2D3748` | Primary dark / backgrounds |
| **Rampart Stone** | `#4A5568` | Secondary dark / borders |
| **Rampart Iron** | `#718096` | Tertiary / muted text |
| **Fortress Gold** | `#D69E2E` | Primary accent / CTAs / highlights |
| **Shield Silver** | `#E2E8F0` | Light backgrounds / cards |

### Semantic Colors

| Name | Hex | Usage |
|------|-----|-------|
| **Accept Green** | `#38A169` | ACCEPT action, success states |
| **Drop Red** | `#E53E3E` | DROP action, errors, deletions |
| **Reject Orange** | `#DD6B20` | REJECT action, warnings |
| **Log Blue** | `#3182CE` | LOG action, info states |
| **Rate Purple** | `#805AD5` | Rate-limit action |

### Dark Mode

| Name | Hex | Usage |
|------|-----|-------|
| **Dark Background** | `#1A202C` | App background |
| **Dark Surface** | `#2D3748` | Cards, panels |
| **Dark Border** | `#4A5568` | Borders, dividers |
| **Dark Text** | `#E2E8F0` | Primary text |
| **Dark Muted** | `#A0AEC0` | Secondary text |

### CSS Variables

```css
:root {
  --rampart-slate: #2D3748;
  --rampart-stone: #4A5568;
  --rampart-iron: #718096;
  --rampart-gold: #D69E2E;
  --rampart-silver: #E2E8F0;
  --rampart-accept: #38A169;
  --rampart-drop: #E53E3E;
  --rampart-reject: #DD6B20;
  --rampart-log: #3182CE;
  --rampart-rate: #805AD5;
}
```

---

## 5. Logo Concepts

### Concept A — Crenellation Wall
A minimalist castle rampart wall (crenellated / merlon pattern) with a code bracket `{ }` integrated into the center merlon. The wall pattern represents defense; the bracket represents policy-as-code.

### Concept B — Shield + Terminal
A shield shape with a terminal cursor `▶` inside. The shield represents protection; the cursor represents CLI-first design.

### Concept C — Lock Glyph
Stylized padlock where the lock body is a brick pattern (wall) and the shackle is a network topology line. Represents network security + infrastructure.

---

## 6. Typography

| Usage | Font | Fallback |
|-------|------|----------|
| **Headings** | Inter Bold | system-ui, sans-serif |
| **Body** | Inter Regular | system-ui, sans-serif |
| **Code / CLI** | JetBrains Mono | monospace |
| **Logo** | Custom / Inter Black | — |

---

## 7. Domain & URLs

| Property | Value |
|----------|-------|
| **Primary domain** | rampart.dev |
| **GitHub org** | github.com/rampartfw |
| **Repository** | github.com/rampartfw/rampart |
| **Documentation** | docs.rampart.dev |
| **Go pkg** | pkg.go.dev/github.com/rampartfw/rampart |
| **Docker Hub** | hub.docker.com/r/rampartfw/rampart |

---

## 8. Nano Banana 2 Prompts

### Prompt 1 — Logo / Icon

```
Create a modern minimalist logo for "Rampart" — a network policy engine.

Design elements:
- A stylized castle rampart wall (crenellated/merlon pattern) in a minimal geometric style
- Integrate code brackets { } into the wall design
- Primary color: Fortress Gold (#D69E2E) for the wall
- Background: Rampart Slate (#2D3748)
- Clean vector style, works at 16px favicon and 512px sizes
- No text in the logo mark
- Military/fortress aesthetic but modern and tech-friendly
```

### Prompt 2 — GitHub Social Preview

```
Create a GitHub social preview image (1280x640) for "Rampart".

Content:
- Title: "Rampart" in bold Inter font
- Subtitle: "Policy-as-Code Firewall. One Binary. Every Backend."
- Left side: Stylized castle rampart wall pattern (crenellated)
- Right side: Terminal showing "rampart plan -f policy.yaml" with colorized diff output
- Background: gradient from Rampart Slate (#2D3748) to dark (#1A202C)
- Accent: Fortress Gold (#D69E2E) for highlights
- Bottom badges: "Go" "nftables" "iptables" "eBPF" "AWS" "GCP" "Azure"
- Clean, professional, developer-focused
```

### Prompt 3 — Architecture Infographic

```
Create a technical architecture infographic for "Rampart" network policy engine.

Layout (vertical flow):
1. TOP: YAML Policy Files (code blocks with syntax highlighting)
2. MIDDLE: Rampart Engine box containing:
   - Parser → Compiler → Conflict Detector → Simulator
   - Snapshot Engine, Audit System, Raft Cluster
3. BOTTOM: Backend boxes arranged horizontally:
   - nftables, iptables, eBPF/XDP, AWS SG, GCP FW, Azure NSG
   - Each with its own icon/color

Color scheme:
- Background: #1A202C (dark)
- Boxes: #2D3748 with #4A5568 borders
- Arrows: #D69E2E (gold)
- Action colors: Accept=#38A169, Drop=#E53E3E
- Text: #E2E8F0
- Style: clean lines, no gradients, tech diagram feel
```

### Prompt 4 — Feature Comparison Table

```
Create a feature comparison infographic: "Rampart vs Manual Firewall Management"

Table format with two columns:
Left (red, old way): "Manual iptables/nftables"
Right (gold/green, new way): "Rampart"

Rows:
- Rules: "iptables -A ... (forgotten)" vs "Version-controlled YAML"
- Audit: "Who changed what? 🤷" vs "Full audit trail + hash chain"
- Rollback: "Hope you remember the old rules" vs "One-click snapshot rollback"
- Multi-host: "Copy-paste + pray" vs "Raft consensus sync"
- Testing: "Apply and hope" vs "Dry-run + packet simulation"
- Time-based: "Set alarm to remember" vs "Auto-expiring rules"
- Cloud: "Different tool per cloud" vs "Unified backend"

Style: dark background (#1A202C), clean icons, Fortress Gold (#D69E2E) accents
```

### Prompt 5 — CLI Screenshot

```
Create a realistic terminal screenshot showing Rampart CLI in action.

Terminal content:
$ rampart plan -f production-web.yaml

Rampart Policy Plan
====================
⚠ 1 warning, 0 errors

WARNING [shadow]: Rule "deny-ssh-all" shadowed by "allow-ssh-bastion"

Plan: 8 rules to add, 2 to remove, 1 to modify.

  + [P10]  allow-ssh-bastion      TCP :22 ← 10.0.1.0/24       ACCEPT
  + [P500] allow-http             TCP :80,:443 ← 0.0.0.0/0    ACCEPT
  ~ [P600] allow-prometheus       TCP :9090,:9100 ← 10.0.10.0/24 ACCEPT
  - [P800] old-temp-debug         TCP :9999 (expired)          REMOVED

Apply? [y/N]: y
✓ Applied 8 rules (6 added, 2 removed, 1 modified) in 15ms

Style: dark terminal theme, green/yellow/red diff colors, monospace font
```

---

## 9. Social Media Content (Turkish / X)

### Tweet 1 — Launch Announcement

```
🏰 Rampart — Network Policy Engine

Firewall kurallarını elle yönetmekten bıktınız mı?

✅ Policy-as-Code (YAML)
✅ Dry-run (Terraform plan gibi)
✅ Instant rollback
✅ Multi-host Raft sync
✅ nftables + iptables + eBPF + Cloud SG

Tek binary. Zero dependency. Pure Go.

github.com/rampartfw/rampart

#golang #devops #security #firewall #opensource
```

### Tweet 2 — Pain Point Thread

```
iptables ile firewall yönetmek:

1. SSH ile sunucuya bağlan
2. iptables -A INPUT -p tcp --dport 22 -j ACCEPT yaz
3. Kural kimin ne zaman eklediğini kimse bilmiyor
4. 6 ay sonra 50 sunucu birbirinden farklı
5. Yanlış kural → SSH kilitleniyor → panic

Rampart ile:
1. YAML yaz
2. rampart plan → diff gör
3. rampart apply → atomic uygula
4. rampart rollback → 1 saniyede geri al

Hangisi?
```

### Tweet 3 — Technical Deep Dive

```
Rampart'ın Conflict Detection Engine'i nasıl çalışıyor?

Rule A: Port 22 ACCEPT from 10.0.0.0/8
Rule B: Port 22 DROP from 10.0.1.0/24

→ B, A'nın subset'i. A zaten 10.0.1.0/24'ü kapsıyor.
→ B asla çalışmaz (shadow conflict).

Rampart bunu compile-time'da yakalıyor. Apply etmeden ÖNCE uyarıyor.

"Terraform plan, but for firewalls."
```

---

## 10. Positioning vs Competitors

| Feature | Rampart | UFW | firewalld | Terraform | Ansible |
|---------|---------|-----|-----------|-----------|---------|
| Policy-as-Code | ✅ YAML | ❌ | ❌ | ✅ HCL | ✅ YAML |
| Dry-run | ✅ | ❌ | ❌ | ✅ | ✅ (--check) |
| Rollback | ✅ Snapshot | ❌ | ❌ | ⚠️ State | ❌ |
| Audit Trail | ✅ Hash chain | ❌ | ❌ | ❌ | ❌ |
| Multi-host Sync | ✅ Raft | ❌ | ❌ | ❌ | Push-based |
| Conflict Detection | ✅ | ❌ | ❌ | ❌ | ❌ |
| Packet Simulation | ✅ | ❌ | ❌ | ❌ | ❌ |
| Time-based Rules | ✅ | ❌ | ⚠️ Rich rules | ❌ | ❌ |
| nftables + iptables | ✅ Both | iptables | nftables | ❌ | ✅ |
| Cloud SG | ✅ | ❌ | ❌ | ✅ | ✅ |
| eBPF/XDP | ✅ | ❌ | ❌ | ❌ | ❌ |
| WebUI | ✅ React | ❌ | Cockpit | ❌ | AWX |
| Single Binary | ✅ Go | ✅ Python | ✅ C/Python | ✅ Go | ❌ Python |
| Zero Dependencies | ✅ | ❌ | ❌ | ✅ | ❌ |
| MCP Server | ✅ | ❌ | ❌ | ❌ | ❌ |
