import { useState, useEffect, createContext, useContext } from 'react'
import { 
  Shield, Cpu, Cloud, Zap, Lock, Menu, X, ChevronRight, Github, 
  Terminal, Activity, Globe, Bot, RefreshCw, Search, CheckCircle2, 
  Sun, Moon, Book, Code2, Layers, ArrowRight, Download, ExternalLink,
  ChevronDown, FileCode, Box, Server, Network, Command, Database, Fingerprint
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const ThemeContext = createContext({ isDark: true, toggle: () => {} })

const ThemeProvider = ({ children }: { children: React.ReactNode }) => {
  const [isDark, setIsDark] = useState(() => {
    const saved = typeof window !== 'undefined' ? localStorage.getItem('rampart-theme') : 'dark'
    return saved === 'dark'
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

const Navbar = () => {
  const { isDark, toggle } = useContext(ThemeContext)
  const [isScrolled, setIsScrolled] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20)
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <nav className={cn(
      "fixed top-0 left-0 right-0 z-[60] transition-all duration-500",
      isScrolled ? "glass border-b border-zinc-200/50 dark:border-white/5 py-4 shadow-2xl" : "bg-transparent py-8"
    )}>
      <div className="container mx-auto px-8 flex items-center justify-between">
        <div className="flex items-center gap-4 group cursor-pointer" onClick={() => window.scrollTo({top: 0, behavior: 'smooth'})}>
          <div className="w-11 h-11 bg-blue-600 rounded-xl flex items-center justify-center shadow-xl shadow-blue-500/30">
            <Shield className="text-white w-6 h-6" />
          </div>
          <span className="text-2xl font-black tracking-tighter dark:text-white text-zinc-950">RAMPART</span>
        </div>

        <div className="hidden md:flex items-center gap-10">
          {['Features', 'Documentation', 'Releases'].map(link => (
            <a key={link} href={`#${link.toLowerCase().replace(' ', '-')}`} className="text-sm font-black text-zinc-500 dark:text-zinc-400 hover:text-blue-600 dark:hover:text-white transition-all uppercase tracking-widest">
              {link}
            </a>
          ))}
          <button onClick={toggle} className="w-10 h-10 rounded-xl flex items-center justify-center bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 hover:border-blue-500/50 transition-all">
            {isDark ? <Sun className="w-4 h-4 text-amber-400" /> : <Moon className="w-4 h-4 text-blue-600" />}
          </button>
          <a href="https://github.com/ersinkoc/Rampart" className="bg-zinc-950 dark:bg-white text-white dark:text-zinc-950 px-6 py-2.5 rounded-xl text-sm font-black transition-all hover:scale-105 shadow-xl">
            GitHub
          </a>
        </div>
        <button className="md:hidden glass p-2 rounded-xl" onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
          {isMobileMenuOpen ? <X /> : <Menu />}
        </button>
      </div>
    </nav>
  )
}

const CodePreview = () => (
  <motion.div initial={{ opacity: 0, y: 40 }} whileInView={{ opacity: 1, y: 0 }} transition={{ duration: 1 }} className="mt-32 max-w-5xl mx-auto glass rounded-[3rem] overflow-hidden shadow-2xl border-zinc-200 dark:border-white/10">
    <div className="bg-zinc-100/80 dark:bg-zinc-900/80 px-8 py-5 border-b border-zinc-200 dark:border-white/10 flex items-center justify-between">
      <div className="flex gap-2.5"><div className="w-3.5 h-3.5 rounded-full bg-rose-500/40" /><div className="w-3.5 h-3.5 rounded-full bg-amber-500/40" /><div className="w-3.5 h-3.5 rounded-full bg-emerald-500/40" /></div>
      <span className="text-[10px] font-black uppercase tracking-[0.4em] font-mono text-zinc-500">policy.yaml</span>
      <div className="w-10" />
    </div>
    <div className="p-10 md:p-16 font-mono text-base md:text-lg text-left bg-white dark:bg-zinc-950/60">
      <pre className="text-zinc-800 dark:text-zinc-300 leading-relaxed overflow-x-auto whitespace-pre">
        <code>{`apiVersion: rampartfw.com/v1\nkind: PolicySet\nmetadata:\n  name: "global-sentinel"\n\npolicies:\n  - name: "threat-mitigation"\n    rules:\n      - name: "block-malicious-dns"\n        action: drop\n        match:\n          appProtocol: dns\n          dns: { query: "attack.io" }`}</code>
      </pre>
    </div>
  </motion.div>
)

const Docs = () => {
  const [tab, setTab] = useState(0)
  const items = [
    { title: 'Install', icon: Terminal, content: 'curl -sSL rampartfw.com/install | sh' },
    { title: 'Raft', icon: Layers, content: 'Distributed consensus with mTLS mühürlü nodes.' },
    { title: 'Sentinel', icon: Bot, content: 'Autonomous threat scoring and cluster blocking.' }
  ]
  return (
    <section id="documentation" className="py-48 bg-zinc-50 dark:bg-zinc-900/20 border-y border-zinc-200 dark:border-white/5">
      <div className="container mx-auto px-8 max-w-6xl flex flex-col lg:flex-row gap-20">
        <div className="lg:w-1/3 space-y-6">
          <h2 className="text-6xl font-black tracking-tighter dark:text-white text-zinc-950 underline decoration-blue-600 decoration-8 underline-offset-4">Docs.</h2>
          {items.map((item, i) => (
            <button key={i} onClick={() => setTab(i)} className={cn("w-full flex items-center gap-6 px-10 py-6 rounded-3xl font-black text-xs uppercase tracking-widest transition-all", tab === i ? "bg-blue-600 text-white shadow-2xl scale-105" : "text-zinc-400 dark:hover:bg-white/5")}>
              <item.icon className="w-5 h-5" /> {item.title}
            </button>
          ))}
        </div>
        <div className="lg:w-2/3 glass rounded-[4rem] p-12 md:p-20 border-zinc-200 dark:border-white/10 bg-white dark:bg-black/40">
          <h3 className="text-4xl font-black mb-8 dark:text-white">{items[tab].title}</h3>
          <p className="text-xl text-zinc-600 dark:text-zinc-400 font-bold mb-10 leading-relaxed">{items[tab].content}</p>
          <div className="p-8 rounded-2xl bg-zinc-950 border border-white/10 font-mono text-emerald-400 shadow-inner">
            <span className="text-zinc-500 mr-2">$</span> {tab === 0 ? 'rampart version' : 'rampart cluster status'}
          </div>
        </div>
      </div>
    </section>
  )
}

const App = () => (
  <ThemeProvider>
    <main className="bg-white dark:bg-zinc-950 min-h-screen text-zinc-900 dark:text-white selection:bg-blue-600/30 font-sans antialiased overflow-x-hidden">
      <Navbar />
      <section className="relative pt-40 pb-48 overflow-hidden">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[150%] h-[1200px] hero-glow pointer-events-none opacity-80" />
        <div className="container mx-auto px-6 relative z-10 text-center max-w-6xl">
          <motion.div initial={{ opacity: 0, y: 30 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.8 }}>
            <span className="inline-flex items-center gap-3 px-5 py-2 rounded-full bg-blue-600/10 border border-blue-500/20 text-blue-600 dark:text-blue-400 text-[10px] font-black uppercase tracking-[0.4em] mb-12 shadow-sm backdrop-blur-md">
              <span className="w-2 h-2 rounded-full bg-blue-500 animate-ping" /> v0.1.0 Initial Beta
            </span>
            <h1 className="text-7xl md:text-[9rem] font-black mb-12 tracking-tighter leading-[0.8] dark:text-white text-zinc-950">Autonomous <br /><span className="text-blue-600 drop-shadow-sm">Defense.</span></h1>
            <p className="text-xl md:text-3xl text-zinc-600 dark:text-zinc-400 mb-16 max-w-4xl mx-auto font-medium tracking-tight">The first policy engine that thinks. Unified management for eBPF, nftables, and every Cloud provider.</p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-8"><button className="w-full sm:w-auto bg-blue-600 hover:bg-blue-500 text-white px-16 py-7 rounded-[2rem] font-black shadow-2xl text-xl">Get Started</button><button className="w-full sm:w-auto glass text-zinc-900 dark:text-white px-16 py-7 rounded-[2rem] font-black flex items-center justify-center gap-4 text-xl group border-2"><Github className="w-6 h-6 text-zinc-500" /> View Source</button></div>
          </motion.div>
          <CodePreview />
        </div>
      </section>

      <section id="features" className="py-48 bg-white dark:bg-zinc-950">
        <div className="container mx-auto px-8 grid md:grid-cols-2 lg:grid-cols-3 gap-12">
          {[
            { icon: Cpu, title: 'eBPF XDP', desc: 'Driver-level filtering performance.' },
            { icon: Bot, title: 'Sentinel', desc: 'Autonomous threat scoring.' },
            { icon: Globe, title: 'Multi-Cloud', desc: 'Unified AWS, GCP, Azure orchestration.' },
            { icon: RefreshCw, title: 'Self-Healing', desc: 'Watchdog drift correction.' },
            { icon: Lock, title: 'mTLS Raft', desc: 'Cluster-wide consistency.' },
            { icon: Search, title: 'L7 DPI', desc: 'Application-aware inspection.' }
          ].map((f, i) => (
            <motion.div key={i} whileHover={{ y: -12 }} className="p-12 rounded-[3.5rem] bg-zinc-50 dark:bg-zinc-900/40 border border-zinc-200 dark:border-white/5 hover:border-blue-500/30 transition-all group">
              <div className="w-16 h-16 rounded-2xl bg-white dark:bg-zinc-800 shadow-xl flex items-center justify-center mb-10 group-hover:bg-blue-600 group-hover:text-white transition-all scale-110"><f.icon className="w-8 h-8" /></div>
              <h3 className="text-3xl font-black mb-4 dark:text-white text-zinc-950 tracking-tighter">{f.title}</h3>
              <p className="text-zinc-500 dark:text-zinc-400 text-lg font-bold leading-relaxed">{f.desc}</p>
            </motion.div>
          ))}
        </div>
      </section>

      <Docs />

      <section id="releases" className="py-48 bg-white dark:bg-zinc-950">
        <div className="container mx-auto px-8 max-w-4xl space-y-12">
          <h2 className="text-6xl font-black tracking-tighter dark:text-white text-zinc-950 flex items-center gap-6"><Box className="w-12 h-12 text-blue-600" /> Releases.</h2>
          <div className="p-12 rounded-[4rem] glass border-zinc-200 dark:border-white/10 flex flex-col md:flex-row md:items-center justify-between gap-12 group transition-all hover:border-blue-500/30">
            <div className="flex gap-10"><div className="w-24 h-24 rounded-[2rem] bg-blue-600 text-white flex items-center justify-center font-black text-3xl shadow-3xl">0.1</div><div className="space-y-4"><h4 className="text-4xl font-black dark:text-white text-zinc-950 tracking-tighter leading-none">v0.1.0 <span className="px-5 py-2 rounded-full bg-emerald-500/10 text-emerald-500 text-[10px] font-black uppercase tracking-[0.3em] ml-4">Stable</span></h4><p className="text-xl text-zinc-500 font-bold">Initial Beta Release • April 2026</p></div></div>
            <a href="https://github.com/ersinkoc/Rampart/releases" className="px-12 py-5 rounded-2xl bg-zinc-950 dark:bg-white text-white dark:text-zinc-950 font-black text-sm flex items-center gap-4 hover:scale-110 active:scale-95 shadow-xl"><Download className="w-5 h-5" /> Assets</a>
          </div>
        </div>
      </section>

      <footer className="py-32 border-t border-zinc-200 dark:border-white/5 bg-zinc-50 dark:bg-black">
        <div className="container mx-auto px-8 flex flex-col md:flex-row items-center justify-between gap-20">
          <div className="flex items-center gap-6 group shrink-0"><div className="w-16 h-16 bg-blue-600 rounded-[1.5rem] flex items-center justify-center shadow-3xl group-hover:rotate-[-12deg] transition-all duration-700"><Shield className="text-white w-9 h-9" /></div><div className="flex flex-col"><span className="text-3xl font-black tracking-tighter dark:text-white text-zinc-950 leading-none">RAMPART</span><span className="text-[10px] font-black uppercase tracking-[0.6em] text-blue-600 mt-2">TACTICAL INTELLIGENCE</span></div></div>
          <p className="text-zinc-500 font-bold text-lg max-w-sm text-center md:text-left leading-relaxed opacity-70">© 2026 Rampart Tactical Intelligence. Secure by design. High-performance by default.</p>
          <div className="flex items-center gap-12 scale-150"><a href="https://github.com/ersinkoc/Rampart" className="text-zinc-400 hover:text-blue-600 transition-all"><Github /></a><a href="#" className="text-zinc-400 hover:text-blue-600 transition-all"><Globe /></a></div>
        </div>
      </footer>
    </main>
  </ThemeProvider>
)

export default App
