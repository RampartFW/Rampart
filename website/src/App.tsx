import { useState, useEffect, createContext, useContext } from 'react'
import { 
  Shield, Cpu, Cloud, Zap, Lock, Menu, X, ChevronRight, GitBranch, 
  Terminal, Globe, Bot, RefreshCw, Search, CheckCircle2, 
  Sun, Moon, Layers, ArrowRight, Download, TerminalSquare,
  Server, Network, Database
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from './lib/utils'

// --- Custom Inline Icons ---
const BrandGithub = ({ className }: { className?: string }) => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}>
    <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 3.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22" />
  </svg>
)

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

// --- Shared UI ---
const Button = ({ className, variant = 'default', size = 'default', ...props }: any) => (
  <button className={cn(
    "inline-flex items-center justify-center whitespace-nowrap rounded-xl text-sm font-bold transition-all active:scale-95 disabled:opacity-50",
    {
      "bg-primary text-primary-foreground hover:opacity-90 shadow-lg shadow-primary/20": variant === 'default',
      "border-2 border-input bg-background hover:bg-accent text-foreground": variant === 'outline',
      "hover:bg-accent text-muted-foreground hover:text-foreground": variant === 'ghost',
      "px-6 py-3": size === 'default',
      "px-10 py-5 text-lg": size === 'lg',
      "w-12 h-12 p-0": size === 'icon'
    },
    className
  )} {...props} />
)

// --- Components ---
const Navbar = () => {
  const { isDark, toggle } = useContext(ThemeContext)
  const [isScrolled, setIsScrolled] = useState(false)
  
  useEffect(() => {
    const h = () => setIsScrolled(window.scrollY > 20)
    window.addEventListener('scroll', h); return () => window.removeEventListener('scroll', h)
  }, [])

  return (
    <nav className={cn(
      "fixed top-0 w-full z-[60] transition-all duration-500",
      isScrolled ? "glass border-b py-3 shadow-xl" : "bg-transparent py-8"
    )}>
      <div className="container mx-auto px-8 flex items-center justify-between">
        <div className="flex items-center gap-3 group cursor-pointer" onClick={() => window.scrollTo({top:0, behavior:'smooth'})}>
          <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-2xl">
            <Shield className="text-white w-6 h-6" />
          </div>
          <span className="text-2xl font-black tracking-tighter text-foreground">RAMPART</span>
        </div>
        
        <div className="hidden md:flex items-center gap-10">
          {['Features', 'Docs', 'Releases'].map(l => (
            <a key={l} href={`#${l.toLowerCase()}`} className="text-sm font-black text-muted-foreground hover:text-primary transition-all tracking-widest">
              {l.toUpperCase()}
            </a>
          ))}
          <div className="flex items-center gap-3">
             <Button variant="outline" size="icon" onClick={toggle}>
                {isDark ? <Sun className="w-5 h-5 text-amber-400" /> : <Moon className="w-5 h-5 text-blue-600" />}
             </Button>
             <a href="https://github.com/ersinkoc/Rampart" target="_blank">
                <Button className="gap-2">
                   <BrandGithub className="w-5 h-5" /> GitHub
                </Button>
             </a>
          </div>
        </div>
      </div>
    </nav>
  )
}

const Hero = () => {
  return (
    <section className="relative pt-40 pb-24 md:pt-64 md:pb-40 overflow-hidden">
      <div className="absolute inset-0 hero-glow pointer-events-none opacity-60" />
      <div className="container mx-auto px-6 relative z-10 text-center max-w-6xl">
        <motion.div initial={{ opacity: 0, y: 30 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.8 }}>
          <span className="inline-flex items-center gap-3 px-5 py-2 rounded-full bg-primary/10 border border-primary/20 text-primary text-[10px] font-black uppercase tracking-[0.4em] mb-12 shadow-sm">
             <span className="w-2 h-2 rounded-full bg-primary animate-ping" /> v0.1.0 Alpha Release
          </span>
          <h1 className="text-6xl md:text-[9.5rem] font-black mb-12 tracking-[-0.06em] leading-[0.85] text-foreground">
             Autonomous <br />
             <span className="text-gradient drop-shadow-sm">Defense.</span>
          </h1>
          <p className="text-xl md:text-3xl text-muted-foreground mb-16 max-w-4xl mx-auto font-medium tracking-tight">
             Rampart is the unified network policy engine that thinks. Zero-trust orchestration for eBPF, nftables, and every Cloud provider.
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-8">
             <Button size="lg" className="px-16 shadow-2xl shadow-primary/30">Get Started Now</Button>
             <Button variant="outline" size="lg" className="px-16 border-2 font-black group">
                <BrandGithub className="w-6 h-6 mr-3 text-muted-foreground group-hover:text-primary transition-colors" /> 
                View Source
             </Button>
          </div>
        </motion.div>
        
        <motion.div initial={{ opacity:0, y:100 }} whileInView={{ opacity:1, y:0 }} viewport={{once:true}} transition={{duration:1}} className="mt-40 max-w-5xl mx-auto glass rounded-[3.5rem] overflow-hidden shadow-2xl border-border">
          <div className="bg-muted/80 px-8 py-5 border-b border-border flex items-center justify-between">
             <div className="flex gap-2.5"><div className="w-3.5 h-3.5 rounded-full bg-rose-500/40" /><div className="w-3.5 h-3.5 rounded-full bg-amber-500/40" /><div className="w-3.5 h-3.5 rounded-full bg-emerald-500/40" /></div>
             <span className="text-[10px] font-black uppercase tracking-[0.3em] font-mono text-muted-foreground">main-policy.yaml</span>
             <div className="w-10" />
          </div>
          <div className="p-10 md:p-16 font-mono text-base md:text-xl text-left bg-black text-[#c9d1d9] overflow-x-auto">
             <pre><code>{`apiVersion: rampartfw.com/v1\nkind: PolicySet\nmetadata:\n  name: "global-fortress"\n\npolicies:\n  - name: "sentinel-protection"\n    rules:\n      - name: "mitigate-dns-flood"\n        action: drop\n        match:\n          appProtocol: dns\n          dns: { query: "attack-source.com" }`}</code></pre>
          </div>
        </motion.div>
      </div>
    </section>
  )
}

const Documentation = () => {
  const [tab, setTab] = useState(0)
  const items = [
    { title: "One-Command Install", icon: Terminal, content: "curl -sSL rampartfw.com/install | sh", desc: "Installs Rampart globally on any Linux or macOS system." },
    { title: "Distributed Raft", icon: Layers, content: "mTLS mühürlü cluster communication.", desc: "Ensures every node in your fleet follows the same policy." },
    { title: "Sentinel Intelligence", icon: Bot, content: "Autonomous risk scoring & mitigation.", desc: "Stops threats in milliseconds without human intervention." }
  ]
  return (
    <section id="docs" className="py-48 bg-muted/30 border-y border-border">
      <div className="container mx-auto px-8 max-w-7xl flex flex-col lg:flex-row gap-24">
        <div className="lg:w-1/3 space-y-8 text-left">
           <h2 className="text-6xl font-black tracking-tighter text-foreground underline decoration-primary decoration-8 underline-offset-8">Docs.</h2>
           <p className="text-xl text-muted-foreground font-bold leading-relaxed">Master the fortress. Orchestrate the global edge.</p>
           <div className="space-y-4 pt-8">
              {items.map((it, i) => (
                <button key={i} onClick={() => setTab(i)} className={cn("w-full flex items-center justify-between px-10 py-6 rounded-3xl font-black text-xs uppercase tracking-widest transition-all", tab === i ? "bg-primary text-primary-foreground shadow-2xl scale-105" : "text-muted-foreground hover:bg-accent/50")}>
                   <div className="flex items-center gap-6"><it.icon className="w-5 h-5" /> {it.title}</div>
                   <ChevronRight className={cn("w-4 h-4 transition-transform", tab === i ? "rotate-90" : "")} />
                </button>
              ))}
           </div>
        </div>
        <div className="lg:w-2/3 glass rounded-[4rem] p-12 md:p-24 border-border relative overflow-hidden bg-background/50">
           <h3 className="text-5xl font-black mb-10 text-foreground tracking-tighter">{items[tab].title}</h3>
           <p className="text-2xl text-muted-foreground font-bold mb-12 leading-relaxed italic">{items[tab].desc}</p>
           <div className="p-10 rounded-[2.5rem] bg-black text-emerald-400 font-mono text-lg shadow-2xl relative">
              <span className="text-zinc-600 mr-4 select-none">$ </span> {items[tab].content}
              <div className="absolute top-4 right-8 text-zinc-800 font-black italic">PRO TERMINAL</div>
           </div>
           <div className="absolute top-0 right-0 w-64 h-64 bg-primary/5 blur-[100px] rounded-full" />
        </div>
      </div>
    </section>
  )
}

const App = () => (
  <ThemeProvider>
    <main className="min-h-screen bg-background">
      <Navbar />
      <Hero />
      <section id="features" className="py-48 container mx-auto px-8">
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-10">
          {[
            { icon: Cpu, title: "eBPF Fast-Path", desc: "XDP driver-level packet filtering performance." },
            { icon: Bot, title: "Self-Learning IPS", desc: "Autonomous threat scoring with mitigation." },
            { icon: Cloud, title: "Unified Cloud", desc: "AWS, GCP, Azure orchestration in one plane." },
            { icon: RefreshCw, title: "Watchdog Engine", desc: "Continuous drift detection and self-healing." },
            { icon: Lock, title: "mTLS Cluster", desc: "Secure distributed consensus via Raft core." },
            { icon: Search, title: "L7 DPI Awareness", desc: "Native DNS and HTTP application intelligence." }
          ].map((f, i) => (
            <motion.div key={i} whileHover={{ y: -12 }} className="premium-card p-12 rounded-[3.5rem] relative overflow-hidden group">
              <div className="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center mb-10 text-primary group-hover:bg-primary group-hover:text-white transition-all scale-110 shadow-xl shadow-primary/10"><f.icon className="w-8 h-8" /></div>
              <h3 className="text-3xl font-black mb-6 text-foreground tracking-tighter">{f.title}</h3>
              <p className="text-lg text-muted-foreground font-bold leading-relaxed opacity-80">{f.desc}</p>
            </motion.div>
          ))}
        </div>
      </section>
      <Documentation />
      <section id="releases" className="py-48 bg-muted/30 border-t border-border">
         <div className="container mx-auto px-4 max-w-4xl text-center">
            <h2 className="text-6xl font-black mb-20 tracking-tighter text-foreground underline decoration-primary decoration-8 underline-offset-8">Releases.</h2>
            <div className="p-12 rounded-[4rem] glass border-border flex flex-col md:flex-row items-center justify-between gap-12 group transition-all hover:scale-105 duration-700 hover:shadow-2xl">
               <div className="flex items-center gap-10">
                  <div className="w-24 h-24 rounded-3xl bg-primary text-white flex items-center justify-center font-black text-4xl shadow-2xl">0.1</div>
                  <div className="text-left"><h4 className="text-4xl font-black text-foreground tracking-tighter leading-none mb-4">v0.1.0 Stable</h4><p className="text-xl text-muted-foreground font-bold italic tracking-tight">Initial Beta Orchestrator • April 2026</p></div>
               </div>
               <a href="https://github.com/ersinkoc/Rampart/releases" className="shrink-0"><Button size="lg" className="px-16 py-6 shadow-2xl shadow-primary/40">Build Assets</Button></a>
            </div>
         </div>
      </section>
      <footer className="py-32 border-t border-border bg-background">
        <div className="container mx-auto px-8 flex flex-col md:flex-row justify-between items-center gap-20">
          <div className="flex items-center gap-4 group shrink-0">
             <div className="w-16 h-16 bg-primary rounded-2xl flex items-center justify-center shadow-3xl group-hover:rotate-12 transition-all duration-700"><Shield className="text-white w-9 h-9" /></div>
             <div className="flex flex-col"><span className="text-3xl font-black tracking-tighter text-foreground leading-none">RAMPART</span><span className="text-[10px] font-black uppercase tracking-[0.6em] text-primary mt-2">AUTONOMOUS DEFENSE</span></div>
          </div>
          <p className="text-muted-foreground font-bold text-lg max-w-sm text-center md:text-left leading-relaxed opacity-60">High-performance network orchestration. Secure by design. Autonomous by default.</p>
          <a href="https://github.com/ersinkoc/Rampart" className="text-muted-foreground hover:text-primary transition-all scale-[2.5] hover:rotate-12"><BrandGithub /></a>
        </div>
      </footer>
    </main>
  </ThemeProvider>
)

export default App
