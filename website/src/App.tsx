import { useState, useEffect, createContext, useContext, useRef } from 'react'
import { 
  Shield, 
  Cpu, 
  Cloud, 
  Zap, 
  Lock, 
  Menu, 
  X, 
  ChevronRight, 
  Github, 
  Terminal,
  Activity,
  Globe,
  Bot,
  RefreshCw,
  Search,
  CheckCircle2,
  Sun,
  Moon,
  Book,
  Code2,
  Layers,
  ArrowRight,
  Download,
  ExternalLink,
  ChevronDown,
  TerminalSquare,
  FileCode,
  Box,
  Server
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

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

// --- Navigation ---
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
    { name: 'Documentation', href: '#docs' },
    { name: 'Architecture', href: '#architecture' },
    { name: 'Releases', href: '#releases' },
  ]

  return (
    <nav className={cn(
      "fixed top-0 left-0 right-0 z-50 transition-all duration-500",
      isScrolled ? "glass border-b border-zinc-200/50 dark:border-white/10 py-3" : "bg-transparent py-6"
    )}>
      <div className="container mx-auto px-6 flex items-center justify-between">
        <div className="flex items-center gap-3 group cursor-pointer" onClick={() => window.scrollTo({top: 0, behavior: 'smooth'})}>
          <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center group-hover:rotate-12 transition-transform shadow-lg shadow-blue-500/20">
            <Shield className="text-white w-6 h-6" />
          </div>
          <span className="text-2xl font-black tracking-tighter dark:text-white text-zinc-900">
            RAMPART
          </span>
        </div>

        <div className="hidden md:flex items-center gap-8">
          {navLinks.map(link => (
            <a key={link.name} href={link.href} className="text-sm font-bold text-zinc-500 dark:text-zinc-400 hover:text-blue-600 dark:hover:text-white transition-colors">
              {link.name}
            </a>
          ))}
          <div className="flex items-center gap-3">
            <button 
              onClick={toggle}
              className="w-9 h-9 rounded-lg flex items-center justify-center bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 hover:bg-zinc-200 dark:hover:bg-white/10 transition-all"
            >
              {isDark ? <Sun className="w-4 h-4 text-amber-400" /> : <Moon className="w-4 h-4 text-blue-600" />}
            </button>
            <a 
              href="https://github.com/ersinkoc/Rampart" 
              className="flex items-center gap-2 bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 px-5 py-2 rounded-xl text-sm font-black transition-all hover:opacity-90 shadow-xl"
            >
              <Github className="w-4 h-4" />
              GitHub
            </a>
          </div>
        </div>

        <button className="md:hidden" onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
          {isMobileMenuOpen ? <X /> : <Menu />}
        </button>
      </div>

      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div 
            initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, y: -20 }}
            className="md:hidden absolute top-full left-0 right-0 glass border-b border-zinc-200 dark:border-white/10 flex flex-col p-8 gap-6 shadow-2xl"
          >
            {navLinks.map(link => (
              <a key={link.name} href={link.href} className="text-xl font-bold dark:text-zinc-300 text-zinc-700" onClick={() => setIsMobileMenuOpen(false)}>
                {link.name}
              </a>
            ))}
            <button onClick={() => { toggle(); setIsMobileMenuOpen(false); }} className="flex items-center justify-between text-xl font-bold dark:text-zinc-300 text-zinc-700">
              Theme
              {isDark ? <Sun className="text-amber-400" /> : <Moon className="text-blue-600" />}
            </button>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  )
}

const Hero = () => {
  return (
    <section className="relative pt-32 pb-20 md:pt-64 md:pb-40 overflow-hidden">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[150%] h-[1200px] hero-gradient pointer-events-none opacity-60" />
      <div className="container mx-auto px-6 relative z-10">
        <div className="max-w-5xl mx-auto text-center">
          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.6 }}>
            <span className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-blue-600/10 border border-blue-500/20 text-blue-600 dark:text-blue-400 text-[10px] font-black uppercase tracking-[0.3em] mb-10 shadow-sm">
              <Zap className="w-3 h-3 fill-current" />
              v0.1.0 Initial Beta is Live
            </span>
            <h1 className="text-6xl md:text-[7.5rem] font-black mb-10 tracking-tighter leading-[0.85] dark:text-white text-zinc-950">
              Autonomous <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-600 via-indigo-500 to-purple-600 dark:from-blue-400 dark:via-blue-600 dark:to-indigo-500">
                Defense Cloud.
              </span>
            </h1>
            <p className="text-xl md:text-2xl text-zinc-600 dark:text-zinc-400 mb-14 leading-relaxed max-w-3xl mx-auto font-medium">
              The unified network policy engine that thinks. Zero-trust orchestration for eBPF, nftables, and every Cloud provider.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-6">
              <button className="w-full sm:w-auto bg-blue-600 hover:bg-blue-500 text-white px-14 py-6 rounded-2xl font-black shadow-2xl shadow-blue-600/40 flex items-center justify-center gap-3 transition-all hover:-translate-y-1 active:scale-95 text-lg">
                Get Started
                <ChevronRight className="w-6 h-6" />
              </button>
              <button className="w-full sm:w-auto glass text-zinc-900 dark:text-white px-14 py-6 rounded-2xl font-black flex items-center justify-center gap-3 transition-all hover:bg-white dark:hover:bg-white/10 text-lg">
                <Github className="w-6 h-6 text-zinc-500" />
                View Source
              </button>
            </div>
          </motion.div>
        </div>

        <motion.div 
          initial={{ opacity: 0, y: 40 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.4, duration: 1 }}
          className="mt-32 max-w-5xl mx-auto glass rounded-[3rem] overflow-hidden shadow-[0_40px_100px_-20px_rgba(0,0,0,0.3)] border-zinc-200 dark:border-white/10"
        >
          <div className="bg-zinc-100/80 dark:bg-zinc-900/80 px-8 py-5 border-b border-zinc-200 dark:border-white/10 flex items-center justify-between">
            <div className="flex gap-2.5">
              <div className="w-3.5 h-3.5 rounded-full bg-rose-500/40" />
              <div className="w-3.5 h-3.5 rounded-full bg-amber-500/40" />
              <div className="w-3.5 h-3.5 rounded-full bg-emerald-500/40" />
            </div>
            <div className="flex items-center gap-2 text-zinc-400">
               <FileCode className="w-4 h-4" />
               <span className="text-[10px] font-black uppercase tracking-[0.4em] font-mono">cluster-policy.yaml</span>
            </div>
            <div className="w-10" />
          </div>
          <div className="p-10 md:p-16 font-mono text-base md:text-lg text-left bg-white dark:bg-zinc-950/60">
            <pre className="text-zinc-800 dark:text-zinc-300 leading-relaxed overflow-x-auto">
              <code>{`apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: "global-ingress"

policies:
  - name: "autonomous-protection"
    rules:
      - name: "mitigate-dns-attacks"
        action: drop
        match:
          appProtocol: dns
          dns: { query: "evil-botnet.com" }
      
      - name: "zero-trust-web"
        action: accept
        match:
          protocol: tcp
          destPorts: [443]`}</code>
            </pre>
          </div>
        </motion.div>
      </div>
    </section>
  )
}

const Documentation = () => {
  const [activeId, setActiveId] = useState('install')
  
  const docs = [
    {
      id: 'install',
      title: "Quick Install",
      icon: TerminalSquare,
      content: (
        <div className="space-y-8">
          <p className="text-lg leading-relaxed text-zinc-600 dark:text-zinc-400">
            Rampart can be installed on any Linux or macOS system with a single command. It automatically detects your architecture and sets up the required kernel bindings.
          </p>
          <div className="bg-zinc-950 p-8 rounded-3xl border border-white/10 font-mono shadow-inner group relative">
             <span className="text-zinc-500 select-none">$ </span>
             <span className="text-emerald-400">curl -sSL rampartfw.com/install | sh</span>
             <button className="absolute right-6 top-1/2 -translate-y-1/2 text-zinc-600 hover:text-white">
                <Code2 className="w-5 h-5" />
             </button>
          </div>
          <div className="grid sm:grid-cols-2 gap-4">
             <div className="p-6 rounded-2xl bg-blue-600/5 border border-blue-500/10">
                <h4 className="font-black text-blue-600 dark:text-blue-400 mb-2 uppercase text-xs tracking-widest">Post-Install</h4>
                <p className="text-sm dark:text-zinc-300">Run <code className="font-bold">rampart version</code> to verify installation.</p>
             </div>
             <div className="p-6 rounded-2xl bg-emerald-600/5 border border-emerald-500/10">
                <h4 className="font-black text-emerald-600 dark:text-emerald-400 mb-2 uppercase text-xs tracking-widest">Self-Healing</h4>
                <p className="text-sm dark:text-zinc-300">Systemd service is automatically enabled on Linux.</p>
             </div>
          </div>
        </div>
      )
    },
    {
      id: 'raft',
      title: "Distributed Core",
      icon: Layers,
      content: (
        <div className="space-y-6">
          <p className="text-lg leading-relaxed text-zinc-600 dark:text-zinc-400">
            Strong consistency is the heart of Rampart. Using an industrial-grade Raft implementation, every node in your cluster is guaranteed to be in sync.
          </p>
          <ul className="space-y-4">
            <li className="flex items-start gap-4">
               <div className="w-6 h-6 rounded-full bg-blue-600/10 flex items-center justify-center mt-1 text-blue-600 font-black text-xs">1</div>
               <div>
                  <h5 className="font-bold dark:text-white">mTLS Everywhere</h5>
                  <p className="text-sm text-zinc-500">All peer-to-peer traffic is encrypted and authenticated with internal CA.</p>
               </div>
            </li>
            <li className="flex items-start gap-4">
               <div className="w-6 h-6 rounded-full bg-blue-600/10 flex items-center justify-center mt-1 text-blue-600 font-black text-xs">2</div>
               <div>
                  <h5 className="font-bold dark:text-white">Autonomous Quorum</h5>
                  <p className="text-sm text-zinc-500">Automatic leader election ensures 100% uptime even if nodes fail.</p>
               </div>
            </li>
          </ul>
        </div>
      )
    },
    {
      id: 'sentinel',
      title: "Autonomous sentinel",
      icon: Bot,
      content: (
        <div className="space-y-6">
          <p className="text-lg leading-relaxed text-zinc-600 dark:text-zinc-400">
            Stop reacting to threats. Rampart's Sentinel module continuously learns from your network traffic.
          </p>
          <div className="p-8 rounded-[2rem] bg-gradient-to-br from-zinc-900 to-zinc-950 border border-white/10 shadow-2xl">
             <div className="flex items-center gap-4 mb-6">
                <div className="w-3 h-3 rounded-full bg-rose-500 animate-pulse" />
                <span className="text-xs font-black uppercase tracking-widest text-rose-500">Alert: Threat Detected</span>
             </div>
             <div className="space-y-2 font-mono text-xs text-zinc-400 mb-8">
                <p>{`> IP 192.168.1.45: Excessive L7 violations`}</p>
                <p>{`> RISK_SCORE: 85/100 (Threshold 70)`}</p>
                <p className="text-emerald-400">{`> ACTION: Cluster-wide block applied.`}</p>
             </div>
             <button className="text-blue-500 text-xs font-black uppercase tracking-widest flex items-center gap-2">
                Configure Sentinel <ArrowRight className="w-3 h-3" />
             </button>
          </div>
        </div>
      )
    }
  ]

  return (
    <section id="docs" className="py-40 bg-zinc-50 dark:bg-zinc-900/20 border-y border-zinc-200 dark:border-white/5">
      <div className="container mx-auto px-6">
        <div className="max-w-6xl mx-auto">
          <div className="flex flex-col lg:flex-row gap-20">
            <div className="lg:w-1/3">
              <h2 className="text-5xl font-black mb-8 tracking-tighter dark:text-white text-zinc-950 underline decoration-blue-600 decoration-8 underline-offset-8">Documentation.</h2>
              <p className="text-zinc-500 dark:text-zinc-400 text-xl mb-12 leading-relaxed font-bold">
                Learn how to deploy and orchestrate your global security policy.
              </p>
              
              <div className="space-y-4">
                {docs.map((doc) => (
                  <button 
                    key={doc.id}
                    onClick={() => setActiveId(doc.id)}
                    className={cn(
                      "w-full flex items-center gap-5 px-8 py-5 rounded-2xl font-black text-sm transition-all text-left",
                      activeId === doc.id 
                        ? "bg-blue-600 text-white shadow-2xl shadow-blue-600/30 scale-105" 
                        : "text-zinc-400 dark:text-zinc-500 hover:bg-zinc-100 dark:hover:bg-white/5"
                    )}
                  >
                    <doc.icon className="w-5 h-5" />
                    {doc.title}
                  </button>
                ))}
              </div>
            </div>

            <div className="lg:w-2/3 glass rounded-[3.5rem] p-10 md:p-20 border-zinc-200 dark:border-white/10 bg-white dark:bg-zinc-950/40 relative">
               <AnimatePresence mode="wait">
                 <motion.div
                   key={activeId}
                   initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: -20 }}
                   transition={{ duration: 0.4 }}
                   className="relative z-10"
                 >
                   <h3 className="text-3xl font-black mb-8 dark:text-white text-zinc-950 tracking-tight">
                     {docs.find(d => d.id === activeId)?.title}
                   </h3>
                   {docs.find(d => d.id === activeId)?.content}
                 </motion.div>
               </AnimatePresence>
               <div className="absolute top-0 right-0 w-64 h-64 bg-blue-600/5 blur-3xl rounded-full" />
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

const Releases = () => {
  const versions = [
    { v: 'v0.1.0', tag: 'Initial Beta', date: 'April 2026', current: true },
    { v: 'v0.1.1', tag: 'Planned', date: 'May 2026', current: false },
  ]

  return (
    <section id="releases" className="py-40">
      <div className="container mx-auto px-6">
        <div className="max-w-4xl mx-auto">
           <div className="flex items-center gap-4 mb-12">
              <Box className="w-10 h-10 text-blue-600" />
              <h2 className="text-5xl font-black tracking-tighter dark:text-white text-zinc-950">Releases.</h2>
           </div>
           
           <div className="space-y-6">
              {versions.map(version => (
                <div key={version.v} className="p-10 rounded-[2.5rem] glass border-zinc-200 dark:border-white/10 flex flex-col md:flex-row md:items-center justify-between gap-8 group">
                   <div className="flex items-center gap-6">
                      <div className={cn(
                        "w-16 h-16 rounded-2xl flex items-center justify-center font-black text-xl shadow-lg",
                        version.current ? "bg-blue-600 text-white" : "bg-zinc-100 dark:bg-zinc-800 text-zinc-400"
                      )}>
                        {version.v.split('.')[1]}.{version.v.split('.')[2]}
                      </div>
                      <div>
                         <h4 className="text-2xl font-black dark:text-white text-zinc-950 flex items-center gap-3">
                           {version.v}
                           {version.current && <span className="px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 text-[10px] uppercase tracking-widest">Stable</span>}
                         </h4>
                         <p className="text-zinc-500 font-bold">{version.tag} • {version.date}</p>
                      </div>
                   </div>
                   
                   <div className="flex gap-4">
                      <a 
                        href={`https://github.com/ersinkoc/Rampart/releases/tag/${version.v}`}
                        className="px-8 py-4 rounded-xl bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 font-black text-sm flex items-center gap-2 hover:opacity-90 transition-all"
                      >
                         <Download className="w-4 h-4" />
                         Assets
                      </a>
                      <button className="p-4 rounded-xl border border-zinc-200 dark:border-white/10 text-zinc-500 hover:text-blue-600 transition-colors">
                         <ExternalLink className="w-5 h-5" />
                      </button>
                   </div>
                </div>
              ))}
           </div>
        </div>
      </div>
    </section>
  )
}

const FeatureGrid = () => {
  const features = [
    { icon: Cpu, title: "eBPF/XDP Offload", desc: "Process millions of packets with driver-level performance and near-zero CPU overhead." },
    { icon: Bot, title: "Adaptive Defense", desc: "Sentinel IPS analyzes DPI signals to autonomously score and mitigate threats cluster-wide." },
    { icon: Cloud, title: "Hybrid-Cloud Sync", desc: "Native orchestration for AWS, GCP, and Azure through a single unified YAML control plane." },
    { icon: RefreshCw, title: "Drift Protection", desc: "Watchdog service continuously audits kernel state and enforces policy compliance automatically." },
    { icon: Globe, title: "Secure Consensus", desc: "Industrial Raft core ensures strong consistency and mTLS security for your entire infrastructure." },
    { icon: Search, title: "L7 Awareness", desc: "Deep Packet Inspection for advanced filtering of DNS queries, HTTP headers, and TLS SNI strings." },
  ]

  return (
    <section id="features" className="py-40 bg-white dark:bg-zinc-950 relative overflow-hidden">
      <div className="container mx-auto px-6">
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-12">
           {features.map((f, i) => (
             <motion.div 
               key={i} whileHover={{ y: -12 }}
               className="p-12 rounded-[3.5rem] bg-zinc-50 dark:bg-zinc-900/40 border border-zinc-200 dark:border-white/5 hover:border-blue-500/30 transition-all group relative overflow-hidden"
             >
                <div className="w-16 h-16 rounded-2xl bg-white dark:bg-zinc-800 shadow-xl flex items-center justify-center mb-10 group-hover:bg-blue-600 group-hover:text-white transition-all scale-110 shadow-blue-600/5">
                   <f.icon className="w-8 h-8 text-blue-600 dark:text-blue-400 group-hover:text-white" />
                </div>
                <h3 className="text-3xl font-black mb-6 dark:text-white text-zinc-950 tracking-tighter leading-none">{f.title}</h3>
                <p className="text-zinc-500 dark:text-zinc-400 text-lg leading-relaxed font-medium">{f.desc}</p>
                <div className="absolute -bottom-8 -right-8 w-24 h-24 bg-blue-600/5 blur-3xl rounded-full opacity-0 group-hover:opacity-100 transition-opacity" />
             </motion.div>
           ))}
        </div>
      </div>
    </section>
  )
}

const App = () => {
  return (
    <ThemeProvider>
      <main className="bg-white dark:bg-zinc-950 min-h-screen text-zinc-900 dark:text-white selection:bg-blue-600/30 font-sans antialiased">
        <Navbar />
        <Hero />
        <FeatureGrid />
        <Documentation />
        <Releases />
        
        <section id="architecture" className="py-40 container mx-auto px-6 text-center">
          <div className="glass p-20 md:p-40 rounded-[5rem] border-zinc-200 dark:border-white/10 relative overflow-hidden bg-zinc-50 dark:bg-zinc-900/60 shadow-[0_50px_100px_-20px_rgba(0,0,0,0.1)]">
            <h2 className="text-6xl md:text-[9rem] font-black mb-12 relative z-10 tracking-[ -0.05em] leading-[0.8] dark:text-white text-zinc-950">
              Future <br />
              <span className="text-blue-600">Proofed.</span>
            </h2>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-10 relative z-10">
              <a href="https://github.com/ersinkoc/Rampart" className="w-full sm:w-auto bg-zinc-950 dark:bg-white text-white dark:text-zinc-950 px-20 py-8 rounded-[2rem] font-black hover:scale-110 transition-all text-2xl shadow-3xl">
                Deploy v0.1.0
              </a>
            </div>
            <div className="absolute top-0 right-0 w-[80%] h-full bg-[radial-gradient(circle_at_top_right,rgba(59,130,246,0.15),transparent)]" />
            <div className="absolute bottom-0 left-0 w-[50%] h-full bg-[radial-gradient(circle_at_bottom_left,rgba(99,102,241,0.05),transparent)]" />
          </div>
        </section>

        <footer className="py-24 border-t border-zinc-200 dark:border-white/5 bg-zinc-50 dark:bg-zinc-950">
          <div className="container mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-16">
            <div className="flex items-center gap-4 group cursor-pointer">
              <div className="w-12 h-12 bg-blue-600 rounded-2xl flex items-center justify-center shadow-2xl shadow-blue-600/20 group-hover:rotate-12 transition-transform">
                <Shield className="text-white w-7 h-7" />
              </div>
              <span className="text-2xl font-black tracking-tighter">RAMPART</span>
            </div>
            <p className="text-zinc-500 font-bold text-sm max-w-sm text-center md:text-left">
              © 2026 Rampart Tactical Intelligence. Secure by design. High-performance by default.
            </p>
            <div className="flex items-center gap-10">
              <a href="https://github.com/ersinkoc/Rampart" className="text-zinc-400 hover:text-blue-600 transition-all scale-150">
                <Github />
              </a>
              <a href="#" className="text-zinc-400 hover:text-blue-600 transition-all scale-150">
                <Globe />
              </a>
            </div>
          </div>
        </footer>
      </main>
    </ThemeProvider>
  )
}

export default App
