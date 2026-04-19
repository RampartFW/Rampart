import { useState, useEffect, createContext, useContext } from 'react'
import { 
  Shield, Cpu, Cloud, Zap, Lock, Menu, X, ChevronRight, GitBranch, 
  Terminal, Globe, Bot, RefreshCw, Search, CheckCircle2, 
  Sun, Moon, Layers, ArrowRight, Download, TerminalSquare,
  Server, Network, Database
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from './lib/utils'

// --- Theme Management ---
const ThemeContext = createContext({ isDark: true, toggle: () => {} })

const ThemeProvider = ({ children }: { children: React.ReactNode }) => {
  const [isDark, setIsDark] = useState(() => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('rampart-theme')
      return saved ? saved === 'dark' : true
    }
    return true
  })

  useEffect(() => {
    const root = window.document.documentElement
    if (isDark) {
      root.classList.add('dark')
      localStorage.setItem('rampart-theme', 'dark')
    } else {
      root.classList.remove('dark')
      localStorage.setItem('rampart-theme', 'light')
    }
  }, [isDark])

  return (
    <ThemeContext.Provider value={{ isDark, toggle: () => setIsDark(!isDark) }}>
      {children}
    </ThemeContext.Provider>
  )
}

// --- Shared Components ---
const Button = ({ className, variant = 'default', size = 'default', ...props }: any) => {
  return (
    <button
      className={cn(
        "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
        {
          "bg-primary text-primary-foreground hover:bg-primary/90 shadow-md": variant === 'default',
          "border border-input bg-background hover:bg-accent hover:text-accent-foreground": variant === 'outline',
          "hover:bg-accent hover:text-accent-foreground": variant === 'ghost',
          "h-10 px-4 py-2": size === 'default',
          "h-12 rounded-md px-8 text-base": size === 'lg',
          "h-10 w-10": size === 'icon',
        },
        className
      )}
      {...props}
    />
  )
}

const Badge = ({ className, variant = "default", ...props }: any) => (
  <div className={cn(
    "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
    {
      "border-transparent bg-primary text-primary-foreground": variant === "default",
      "border-transparent bg-secondary text-secondary-foreground": variant === "secondary",
      "text-foreground": variant === "outline",
    },
    className
  )} {...props} />
)

// --- Sections ---
const Navbar = () => {
  const { isDark, toggle } = useContext(ThemeContext)
  const [isScrolled, setIsScrolled] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20)
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  const navLinks = [
    { name: 'Features', href: '#features' },
    { name: 'Docs', href: '#docs' },
    { name: 'Releases', href: '#releases' },
  ]

  return (
    <header className={cn(
      "fixed top-0 w-full z-50 transition-all duration-300 border-b",
      isScrolled ? "bg-background/80 backdrop-blur-md border-border shadow-sm" : "bg-transparent border-transparent"
    )}>
      <div className="container mx-auto px-4 md:px-8 h-16 flex items-center justify-between">
        <div className="flex items-center gap-2 cursor-pointer" onClick={() => window.scrollTo({top: 0, behavior: 'smooth'})}>
          <Shield className="w-6 h-6 text-primary" />
          <span className="font-bold tracking-tight text-lg">Rampart</span>
        </div>

        <nav className="hidden md:flex items-center gap-6">
          {navLinks.map(link => (
            <a key={link.name} href={link.href} className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
              {link.name}
            </a>
          ))}
        </nav>

        <div className="hidden md:flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={toggle} className="text-muted-foreground">
            {isDark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
          </Button>
          <a href="https://github.com/ersinkoc/Rampart" target="_blank" rel="noreferrer">
            <Button variant="outline" size="sm" className="gap-2">
              <GitBranch className="w-4 h-4" /> GitHub
            </Button>
          </a>
        </div>

        <Button variant="ghost" size="icon" className="md:hidden" onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
          {isMobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
        </Button>
      </div>

      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div 
            initial={{ opacity: 0, y: -10 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, y: -10 }}
            className="md:hidden absolute top-16 left-0 w-full bg-background border-b border-border shadow-lg p-4 flex flex-col gap-4"
          >
            {navLinks.map(link => (
              <a key={link.name} href={link.href} className="text-sm font-medium p-2 hover:bg-accent rounded-md" onClick={() => setIsMobileMenuOpen(false)}>
                {link.name}
              </a>
            ))}
            <div className="h-px bg-border my-2" />
            <div className="flex items-center justify-between p-2">
              <span className="text-sm font-medium">Theme</span>
              <Button variant="ghost" size="icon" onClick={toggle}>
                {isDark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
              </Button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  )
}

const Hero = () => {
  const { isDark } = useContext(ThemeContext)
  return (
    <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden">
      <div className={cn("absolute inset-0 pointer-events-none opacity-50", isDark ? "glow-hero-dark" : "glow-hero-light")} />
      
      <div className="container mx-auto px-4 md:px-8 relative z-10 text-center">
        <div className="max-w-4xl mx-auto">
          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }}>
            <Badge variant="secondary" className="mb-6 py-1 px-3 gap-1 border-primary/20 text-primary">
              <Zap className="w-3 h-3" /> v0.1.0 Initial Beta is Live
            </Badge>
            
            <h1 className="text-5xl md:text-7xl font-bold tracking-tight mb-6 leading-tight">
              The Sovereign <br className="hidden md:block" />
              <span className="text-primary">Defense Engine.</span>
            </h1>
            
            <p className="text-lg md:text-xl text-muted-foreground mb-10 max-w-2xl mx-auto leading-relaxed">
              Stop managing firewall rules by hand. Rampart is a unified, autonomous policy engine for Linux eBPF, nftables, and Cloud Security Groups.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <a href="#docs">
                <Button size="lg" className="w-full sm:w-auto gap-2 font-semibold">
                  Get Started <ChevronRight className="w-4 h-4" />
                </Button>
              </a>
              <a href="https://github.com/ersinkoc/Rampart" target="_blank" rel="noreferrer">
                <Button variant="outline" size="lg" className="w-full sm:w-auto gap-2 font-semibold">
                  <GitBranch className="w-4 h-4" /> View Source
                </Button>
              </a>
            </div>
          </motion.div>
        </div>

        <motion.div 
          initial={{ opacity: 0, y: 40 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2, duration: 0.7 }}
          className="mt-20 max-w-5xl mx-auto rounded-2xl border border-border bg-card shadow-2xl overflow-hidden"
        >
          <div className="flex items-center gap-2 px-4 py-3 border-b border-border bg-muted/50">
            <div className="flex gap-1.5">
              <div className="w-3 h-3 rounded-full bg-rose-500/80" />
              <div className="w-3 h-3 rounded-full bg-amber-500/80" />
              <div className="w-3 h-3 rounded-full bg-emerald-500/80" />
            </div>
            <span className="ml-2 text-xs font-mono text-muted-foreground">policy.yaml</span>
          </div>
          <div className="p-6 md:p-8 font-mono text-sm md:text-base bg-[#0d1117] text-[#c9d1d9] text-left overflow-x-auto">
            <pre className="leading-loose">
              <code>
<span className="text-[#ff7b72]">apiVersion</span>: <span className="text-[#a5d6ff]">rampartfw.com/v1</span>{`\n`}
<span className="text-[#ff7b72]">kind</span>: <span className="text-[#a5d6ff]">PolicySet</span>{`\n`}
<span className="text-[#ff7b72]">metadata</span>:{`\n`}
{`  `}<span className="text-[#79c0ff]">name</span>: <span className="text-[#a5d6ff]">"global-sentinel"</span>{`\n\n`}
<span className="text-[#8b949e]"># Autonomous Security Logic</span>{`\n`}
<span className="text-[#ff7b72]">policies</span>:{`\n`}
{`  - `}<span className="text-[#79c0ff]">name</span>: <span className="text-[#a5d6ff]">"zero-trust-edge"</span>{`\n`}
{`    `}<span className="text-[#79c0ff]">rules</span>:{`\n`}
{`      - `}<span className="text-[#79c0ff]">name</span>: <span className="text-[#a5d6ff]">"mitigate-l7-threats"</span>{`\n`}
{`        `}<span className="text-[#79c0ff]">action</span>: <span className="text-[#ff7b72] font-bold">drop</span>{`\n`}
{`        `}<span className="text-[#79c0ff]">match</span>:{`\n`}
{`          `}<span className="text-[#79c0ff]">appProtocol</span>: <span className="text-[#a5d6ff]">dns</span>{`\n`}
{`          `}<span className="text-[#79c0ff]">dns</span>: <span className="text-[#c9d1d9]">{`{ query: "attack-source.com" }`}</span>
              </code>
            </pre>
          </div>
        </motion.div>
      </div>
    </section>
  )
}

const Features = () => {
  const items = [
    { icon: Cpu, title: "eBPF/XDP Fast-Path", desc: "Native kernel execution with microsecond latency. Offload massive rule sets directly to the network driver." },
    { icon: Bot, title: "Autonomous Sentinel", desc: "Beyond traditional IPS. Real-time threat scoring and automated cluster-wide mitigation powered by DPI signals." },
    { icon: Cloud, title: "Unified Cloud Plane", desc: "Single control plane for AWS, GCP, and Azure. Synchronize security groups globally with one YAML manifest." },
    { icon: RefreshCw, title: "Self-Healing Watchdog", desc: "Continuous drift detection. Rampart automatically audits and corrects firewall state to ensure 100% compliance." },
    { icon: Lock, title: "mTLS Raft Cluster", desc: "Industrial Raft core ensures strong consistency and encrypted peer-to-peer security for your entire infrastructure." },
    { icon: Search, title: "Layer-7 DPI Intelligence", desc: "Deep Packet Inspection for DNS, HTTP, and TLS SNI. Stop application-level threats before they enter your network." },
  ]

  return (
    <section id="features" className="py-24 bg-muted/30 border-y border-border">
      <div className="container mx-auto px-4 md:px-8">
        <div className="text-center max-w-2xl mx-auto mb-16">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">Built for Resilience</h2>
          <p className="text-muted-foreground text-lg">Traditional firewalls are static and complex. Rampart is dynamic, distributed, and aware.</p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
           {items.map((f, i) => (
             <div key={i} className="p-6 rounded-2xl bg-card border border-border shadow-sm hover:shadow-md transition-shadow">
                <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4 text-primary">
                   <f.icon className="w-6 h-6" />
                </div>
                <h3 className="text-lg font-semibold mb-2">{f.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{f.desc}</p>
             </div>
           ))}
        </div>
      </div>
    </section>
  )
}

const Docs = () => {
  const [activeTab, setActiveTab] = useState(0)
  
  const docs = [
    {
      title: "Quick Start",
      icon: TerminalSquare,
      content: (
        <div className="space-y-6 animate-in fade-in duration-500">
          <div>
             <h3 className="text-2xl font-bold tracking-tight mb-2">Deployment</h3>
             <p className="text-muted-foreground">Rampart is a single static binary. Install the controller and agent globally in seconds.</p>
          </div>
          <div className="code-block">
             <span className="text-muted-foreground select-none">$ </span>
             <span className="text-emerald-400">curl -sSL rampartfw.com/install | sh</span>
          </div>
          <div className="grid sm:grid-cols-2 gap-4 mt-6">
             <div className="p-5 rounded-xl border border-border bg-muted/50">
                <CheckCircle2 className="w-5 h-5 text-primary mb-2" />
                <h4 className="font-semibold text-sm mb-1">Auto Discovery</h4>
                <p className="text-xs text-muted-foreground">Detects OS, Arch, and kernel features (XDP) automatically.</p>
             </div>
             <div className="p-5 rounded-xl border border-border bg-muted/50">
                <CheckCircle2 className="w-5 h-5 text-primary mb-2" />
                <h4 className="font-semibold text-sm mb-1">Systemd Ready</h4>
                <p className="text-xs text-muted-foreground">Hooks into init systems with pre-configured self-healing services.</p>
             </div>
          </div>
        </div>
      )
    },
    {
      title: "Architecture",
      icon: Network,
      content: (
        <div className="space-y-6 animate-in fade-in duration-500">
           <div>
             <h3 className="text-2xl font-bold tracking-tight mb-2">Distributed Control Plane</h3>
             <p className="text-muted-foreground">Rampart is designed for zero-single-point-of-failure clusters.</p>
          </div>
          <div className="p-6 rounded-xl bg-card border border-border space-y-4">
             <div className="flex items-start gap-4">
                <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0 text-primary font-bold text-sm">1</div>
                <div>
                   <h5 className="font-semibold text-sm">mTLS Everywhere</h5>
                   <p className="text-sm text-muted-foreground mt-1">All peer-to-peer traffic is encrypted and authenticated via internal CA.</p>
                </div>
             </div>
             <div className="flex items-start gap-4">
                <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0 text-primary font-bold text-sm">2</div>
                <div>
                   <h5 className="font-semibold text-sm">Autonomous Quorum</h5>
                   <p className="text-sm text-muted-foreground mt-1">Automatic leader election via Raft ensures 100% uptime during partitions.</p>
                </div>
             </div>
          </div>
        </div>
      )
    },
    {
      title: "Sentinel IPS",
      icon: Bot,
      content: (
        <div className="space-y-6 animate-in fade-in duration-500">
           <div>
             <h3 className="text-2xl font-bold tracking-tight mb-2">Otonom Threat Response</h3>
             <p className="text-muted-foreground">Stop threats in milliseconds without human intervention.</p>
          </div>
          <div className="code-block border-rose-500/20 bg-[#0d1117]">
             <div className="flex items-center gap-2 mb-4">
                <div className="w-2 h-2 rounded-full bg-rose-500 animate-pulse" />
                <span className="text-xs font-bold uppercase tracking-widest text-rose-500">Sentinel Alert</span>
             </div>
             <p className="text-muted-foreground text-xs mb-2"># Real-time traffic analysis result</p>
             <p className="text-sm text-zinc-300">{`> SCORING: IP 192.168.1.45 -> RISK 85 (Target 70)`}</p>
             <p className="text-sm text-rose-400 font-semibold mt-1">{`> ACTION: Cluster-wide block applied via Raft.`}</p>
          </div>
        </div>
      )
    },
  ]

  return (
    <section id="docs" className="py-24">
      <div className="container mx-auto px-4 md:px-8">
        <div className="flex flex-col lg:flex-row gap-12">
          <div className="lg:w-1/3">
            <h2 className="text-3xl font-bold mb-4 tracking-tight">Documentation</h2>
            <p className="text-muted-foreground mb-8">Learn how to deploy and orchestrate your global security policy.</p>
            
            <div className="space-y-2 flex flex-col">
              {docs.map((doc, idx) => (
                <button 
                  key={idx}
                  onClick={() => setActiveTab(idx)}
                  className={cn(
                    "flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-all text-left",
                    activeTab === idx 
                      ? "bg-primary text-primary-foreground shadow-sm" 
                      : "text-muted-foreground hover:bg-muted hover:text-foreground"
                  )}
                >
                  <doc.icon className="w-4 h-4" />
                  {doc.title}
                  {activeTab === idx && <ArrowRight className="w-4 h-4 ml-auto opacity-50" />}
                </button>
              ))}
            </div>
          </div>

          <div className="lg:w-2/3">
             <div className="rounded-2xl border border-border bg-card p-6 md:p-10 shadow-sm min-h-[400px]">
               <AnimatePresence mode="wait">
                 <motion.div
                   key={activeTab}
                   initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, y: -10 }}
                   transition={{ duration: 0.2 }}
                 >
                   {docs[activeTab].content}
                 </motion.div>
               </AnimatePresence>
             </div>
          </div>
        </div>
      </div>
    </section>
  )
}

const Releases = () => {
  return (
    <section id="releases" className="py-24 bg-muted/30 border-t border-border">
      <div className="container mx-auto px-4 md:px-8">
        <div className="max-w-4xl mx-auto">
           <h2 className="text-3xl font-bold tracking-tight mb-8">Releases</h2>
           
           <div className="p-6 md:p-8 rounded-2xl bg-card border border-border flex flex-col md:flex-row md:items-center justify-between gap-6 shadow-sm">
               <div className="flex items-center gap-6">
                  <div className="w-16 h-16 rounded-xl bg-primary/10 text-primary flex items-center justify-center font-bold text-2xl">
                    0.1
                  </div>
                  <div>
                     <h4 className="text-2xl font-bold flex items-center gap-3">
                       v0.1.0
                       <Badge>Stable</Badge>
                     </h4>
                     <p className="text-sm text-muted-foreground mt-1">Initial Beta Release • April 2026</p>
                  </div>
               </div>
               
               <a href="https://github.com/ersinkoc/Rampart/releases/tag/v0.1.0">
                 <Button className="gap-2 w-full md:w-auto">
                   <Download className="w-4 h-4" /> Get Assets
                 </Button>
               </a>
           </div>
        </div>
      </div>
    </section>
  )
}

const Footer = () => (
  <footer className="py-12 border-t border-border bg-background">
    <div className="container mx-auto px-4 md:px-8 flex flex-col md:flex-row items-center justify-between gap-6">
      <div className="flex items-center gap-2">
        <Shield className="text-primary w-5 h-5" />
        <span className="font-bold tracking-tight">Rampart</span>
      </div>
      <p className="text-muted-foreground text-sm text-center">
        © 2026 Rampart Tactical Intelligence. Licensed under Apache 2.0.
      </p>
      <div className="flex items-center gap-4">
        <a href="https://github.com/ersinkoc/Rampart" className="text-muted-foreground hover:text-foreground transition-colors">
          <GitBranch className="w-5 h-5" />
        </a>
      </div>
    </div>
  </footer>
)

const App = () => {
  return (
    <ThemeProvider>
      <main className="min-h-screen">
        <Navbar />
        <Hero />
        <Features />
        <Docs />
        <Releases />
        <Footer />
      </main>
    </ThemeProvider>
  )
}

export default App
