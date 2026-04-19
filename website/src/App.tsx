import { useState, useEffect, createContext, useContext } from 'react'
import { 
  Shield, Cpu, Cloud, Zap, Lock, Menu, X, ChevronRight, 
  Terminal, Globe, Bot, RefreshCw, Search, CheckCircle2, 
  Sun, Moon, Layers, ArrowRight, Download, TerminalSquare,
  Server, Network, Database
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from './lib/utils'

const GithubIcon = ({ className }: { className?: string }) => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}>
    <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 3.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22" />
  </svg>
)

const ThemeContext = createContext({ isDark: true, toggle: () => {} })

const ThemeProvider = ({ children }: { children: React.ReactNode }) => {
  const [isDark, setIsDark] = useState(() => {
    const saved = typeof window !== 'undefined' ? localStorage.getItem('rampart-theme') : 'dark'
    return saved !== 'light'
  })
  useEffect(() => {
    const root = window.document.documentElement
    if (isDark) { root.classList.add('dark'); localStorage.setItem('rampart-theme', 'dark') }
    else { root.classList.remove('dark'); localStorage.setItem('rampart-theme', 'light') }
  }, [isDark])
  return <ThemeContext.Provider value={{ isDark, toggle: () => setIsDark(!isDark) }}>{children}</ThemeContext.Provider>
}

const Button = ({ className, variant = 'default', size = 'default', ...props }: any) => (
  <button className={cn("inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors disabled:opacity-50", {
    "bg-primary text-primary-foreground hover:bg-primary/90 shadow-md": variant === 'default',
    "border border-input bg-background hover:bg-accent": variant === 'outline',
    "hover:bg-accent": variant === 'ghost',
    "h-10 px-4 py-2": size === 'default',
    "h-12 px-8 text-base": size === 'lg',
    "h-10 w-10": size === 'icon'
  }, className)} {...props} />
)

const Badge = ({ className, ...props }: any) => (
  <div className={cn("inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold bg-primary text-primary-foreground", className)} {...props} />
)

const Navbar = () => {
  const { isDark, toggle } = useContext(ThemeContext)
  const [scrolled, setScrolled] = useState(false)
  useEffect(() => {
    const h = () => setScrolled(window.scrollY > 20)
    window.addEventListener('scroll', h); return () => window.removeEventListener('scroll', h)
  }, [])
  return (
    <header className={cn("fixed top-0 w-full z-50 transition-all border-b", scrolled ? "bg-background/80 backdrop-blur-md border-border shadow-sm" : "bg-transparent border-transparent")}>
      <div className="container mx-auto px-4 md:px-8 h-16 flex items-center justify-between">
        <div className="flex items-center gap-2 cursor-pointer" onClick={() => window.scrollTo({top: 0, behavior: 'smooth'})}>
          <Shield className="w-6 h-6 text-primary" /><span className="font-bold tracking-tight text-lg text-foreground">Rampart</span>
        </div>
        <nav className="hidden md:flex items-center gap-6">
          {['Features', 'Docs', 'Releases'].map(l => <a key={l} href={`#${l.toLowerCase()}`} className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">{l}</a>)}
        </nav>
        <div className="hidden md:flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={toggle} className="text-muted-foreground">{isDark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}</Button>
          <a href="https://github.com/ersinkoc/Rampart" target="_blank" rel="noreferrer"><Button variant="outline" size="sm" className="gap-2 font-bold"><GithubIcon className="w-4 h-4" /> GitHub</Button></a>
        </div>
      </div>
    </header>
  )
}

const Hero = () => {
  const { isDark } = useContext(ThemeContext)
  return (
    <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden text-center text-foreground">
      <div className={cn("absolute inset-0 pointer-events-none opacity-50", isDark ? "glow-hero-dark" : "glow-hero-light")} />
      <div className="container mx-auto px-4 relative z-10">
        <Badge className="mb-6 py-1 px-3 gap-1 bg-secondary text-primary border-primary/20"><Zap className="w-3 h-3" /> v0.1.0 Initial Beta</Badge>
        <h1 className="text-5xl md:text-7xl font-bold tracking-tight mb-6">The Sovereign <span className="text-primary">Defense Engine.</span></h1>
        <p className="text-lg md:text-xl text-muted-foreground mb-10 max-w-2xl mx-auto">Autonomous policy engine for Linux eBPF, nftables, and Cloud Security Groups.</p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <a href="#docs"><Button size="lg" className="font-bold">Get Started <ChevronRight className="w-4 h-4" /></Button></a>
          <a href="https://github.com/ersinkoc/Rampart" target="_blank" rel="noreferrer"><Button variant="outline" size="lg" className="gap-2 font-bold"><GithubIcon className="w-4 h-4" /> View Source</Button></a>
        </div>
        <motion.div initial={{ opacity: 0, y: 40 }} animate={{ opacity: 1, y: 0 }} className="mt-20 max-w-4xl mx-auto rounded-2xl border border-border bg-card shadow-2xl overflow-hidden text-left">
          <div className="px-4 py-3 border-b border-border bg-muted/50 flex gap-1.5"><div className="w-3 h-3 rounded-full bg-rose-500/80" /><div className="w-3 h-3 rounded-full bg-amber-500/80" /><div className="w-3 h-3 rounded-full bg-emerald-500/80" /></div>
          <div className="p-6 font-mono text-sm bg-[#0d1117] text-[#c9d1d9] overflow-x-auto">
            <pre><code>{`apiVersion: rampartfw.com/v1\nkind: PolicySet\nmetadata:\n  name: "global-sentinel"\n\npolicies:\n  - name: "zero-trust-edge"\n    rules:\n      - name: "mitigate-dns-threats"\n        action: drop\n        match:\n          appProtocol: dns\n          dns: { query: "attack.com" }`}</code></pre>
          </div>
        </motion.div>
      </div>
    </section>
  )
}

const Features = () => (
  <section id="features" className="py-24 bg-muted/30 border-y border-border">
    <div className="container mx-auto px-4 text-center">
      <h2 className="text-3xl md:text-4xl font-bold mb-16">Built for Resilience</h2>
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        {[
          { icon: Cpu, title: "eBPF/XDP", desc: "Driver-level filtering performance." },
          { icon: Bot, title: "Sentinel IPS", desc: "Autonomous threat detection." },
          { icon: Cloud, title: "Unified Cloud", desc: "AWS, GCP, Azure orchestration." },
          { icon: RefreshCw, title: "Self-Healing", desc: "Continuous drift detection." },
          { icon: Lock, title: "mTLS Raft", desc: "Strong cluster consistency." },
          { icon: Search, title: "L7 DPI", desc: "Application-aware security." }
        ].map((f, i) => (
          <div key={i} className="p-6 rounded-2xl bg-card border border-border text-left shadow-sm">
            <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4 text-primary"><f.icon className="w-6 h-6" /></div>
            <h3 className="text-lg font-semibold mb-2">{f.title}</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">{f.desc}</p>
          </div>
        ))}
      </div>
    </div>
  </section>
)

const Docs = () => {
  const [tab, setTab] = useState(0)
  const items = [
    { title: 'Install', icon: TerminalSquare, content: 'curl -sSL rampartfw.com/install | sh' },
    { title: 'Architecture', icon: Network, content: 'Distributed control plane via Raft consensus.' },
    { icon: Bot, title: 'Sentinel', content: 'Autonomous risk scoring and cluster-wide block.' }
  ]
  return (
    <section id="docs" className="py-24">
      <div className="container mx-auto px-4 flex flex-col lg:flex-row gap-12">
        <div className="lg:w-1/3">
          <h2 className="text-3xl font-bold mb-8">Documentation</h2>
          <div className="space-y-2 flex flex-col">
            {items.map((it, i) => (
              <button key={i} onClick={() => setTab(i)} className={cn("flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-bold transition-all text-left", tab === i ? "bg-primary text-primary-foreground shadow-sm" : "text-muted-foreground hover:bg-muted")}>
                <it.icon className="w-4 h-4" /> {it.title}
              </button>
            ))}
          </div>
        </div>
        <div className="lg:w-2/3 p-10 rounded-2xl border border-border bg-card shadow-sm min-h-[300px]">
          <h3 className="text-2xl font-bold mb-4">{items[tab].title}</h3>
          <p className="text-muted-foreground mb-6 font-medium">{items[tab].content}</p>
          <div className="p-4 bg-zinc-950 text-emerald-400 rounded-xl font-mono text-sm shadow-inner"><span className="text-muted-foreground">$</span> rampart version</div>
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
      <Features />
      <Docs />
      <section id="releases" className="py-24 bg-muted/30 border-t border-border">
        <div className="container mx-auto px-4 max-w-4xl">
          <h2 className="text-3xl font-bold mb-8">Releases</h2>
          <div className="p-8 rounded-2xl bg-card border border-border flex items-center justify-between shadow-sm">
            <div className="flex items-center gap-6"><div className="w-16 h-16 rounded-xl bg-primary/10 text-primary flex items-center justify-center font-bold text-2xl text-foreground">0.1</div>
            <div><h4 className="text-2xl font-bold">v0.1.0</h4><p className="text-sm text-muted-foreground mt-1">Stable Release • April 2026</p></div></div>
            <a href="https://github.com/ersinkoc/Rampart/releases"><Button className="gap-2 font-bold"><Download className="w-4 h-4" /> Assets</Button></a>
          </div>
        </div>
      </section>
      <footer className="py-12 border-t border-border bg-background">
        <div className="container mx-auto px-4 flex justify-between items-center">
          <div className="flex items-center gap-2"><Shield className="text-primary w-5 h-5" /><span className="font-bold tracking-tight">Rampart</span></div>
          <p className="text-muted-foreground text-sm">© 2026 Rampart Tactical Intelligence.</p>
          <a href="https://github.com/ersinkoc/Rampart" className="text-muted-foreground hover:text-foreground"><GithubIcon className="w-5 h-5" /></a>
        </div>
      </footer>
    </main>
  </ThemeProvider>
)

export default App
