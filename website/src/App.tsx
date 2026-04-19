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
  Server,
  Network,
  Command,
  Database,
  Fingerprint
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

// --- Modern Glass Navigation ---
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
    { name: 'Capabilities', href: '#features' },
    { name: 'Documentation', href: '#docs' },
    { name: 'Architecture', href: '#architecture' },
    { name: 'Release Notes', href: '#releases' },
  ]

  return (
    <nav className={cn(
      "fixed top-0 left-0 right-0 z-[60] transition-all duration-700",
      isScrolled ? "glass border-b border-zinc-200/50 dark:border-white/5 py-4 shadow-2xl shadow-black/10" : "bg-transparent py-8"
    )}>
      <div className="container mx-auto px-8 flex items-center justify-between">
        <div className="flex items-center gap-4 group cursor-pointer" onClick={() => window.scrollTo({top: 0, behavior: 'smooth'})}>
          <div className="w-11 h-11 bg-blue-600 rounded-[0.8rem] flex items-center justify-center group-hover:scale-110 group-hover:rotate-6 transition-all duration-500 shadow-xl shadow-blue-500/30">
            <Shield className="text-white w-6 h-6" />
          </div>
          <div className="flex flex-col">
            <span className="text-2xl font-black tracking-[-0.05em] dark:text-white text-zinc-950 leading-none">
              RAMPART
            </span>
            <span className="text-[10px] font-black uppercase tracking-[0.4em] text-blue-600 mt-1">SENTINEL</span>
          </div>
        </div>

        <div className="hidden md:flex items-center gap-12">
          {navLinks.map(link => (
            <a key={link.name} href={link.href} className="text-sm font-black text-zinc-500 dark:text-zinc-500 hover:text-blue-600 dark:hover:text-white transition-all uppercase tracking-widest">
              {link.name}
            </a>
          ))}
          <div className="flex items-center gap-4">
            <button 
              onClick={toggle}
              className="w-10 h-10 rounded-xl flex items-center justify-center bg-zinc-100/50 dark:bg-white/5 border border-zinc-200 dark:border-white/10 hover:border-blue-500/50 transition-all"
            >
              {isDark ? <Sun className="w-4 h-4 text-amber-400" /> : <Moon className="w-4 h-4 text-blue-600" />}
            </button>
            <a 
              href="https://github.com/ersinkoc/Rampart" 
              className="flex items-center gap-2 bg-zinc-950 dark:bg-white text-white dark:text-zinc-950 px-6 py-2.5 rounded-[0.75rem] text-sm font-black transition-all hover:scale-105 shadow-xl shadow-black/20 active:scale-95"
            >
              <Github className="w-4 h-4" />
              v0.1.0
            </a>
          </div>
        </div>

        <button className="md:hidden p-2 rounded-xl glass" onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
          {isMobileMenuOpen ? <X /> : <Menu />}
        </button>
      </div>

      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div 
            initial={{ opacity: 0, height: 0 }} animate={{ opacity: 1, height: 'auto' }} exit={{ opacity: 0, height: 0 }}
            className="md:hidden absolute top-full left-0 right-0 glass border-b border-zinc-200 dark:border-white/10 flex flex-col p-8 gap-6 shadow-2xl overflow-hidden"
          >
            {navLinks.map(link => (
              <a key={link.name} href={link.href} className="text-2xl font-black dark:text-zinc-300 text-zinc-800 tracking-tighter" onClick={() => setIsMobileMenuOpen(false)}>
                {link.name}
              </a>
            ))}
            <hr className="border-zinc-200 dark:border-white/10" />
            <button onClick={() => { toggle(); setIsMobileMenuOpen(false); }} className="flex items-center justify-between text-xl font-bold dark:text-zinc-300 text-zinc-800">
              Theme Mode
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
    <section className="relative pt-40 pb-20 md:pt-64 md:pb-48 overflow-hidden glow-mesh">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[150%] h-[1200px] hero-glow pointer-events-none opacity-80" />
      
      <div className="container mx-auto px-8 relative z-10">
        <div className="max-w-6xl mx-auto text-center">
          <motion.div initial={{ opacity: 0, y: 30 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.8, ease: "easeOut" }}>
            <div className="inline-flex items-center gap-3 px-5 py-2 rounded-full bg-blue-600/10 border border-blue-500/20 text-blue-600 dark:text-blue-400 text-[11px] font-black uppercase tracking-[0.4em] mb-12 shadow-sm backdrop-blur-md">
              <span className="w-2 h-2 rounded-full bg-blue-500 animate-ping" />
              The Sovereign Defense Engine
            </div>
            
            <h1 className="text-7xl md:text-[9.5rem] font-black mb-12 tracking-[-0.07em] leading-[0.8] dark:text-white text-zinc-950">
              Stop Reacting. <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-br from-blue-600 via-indigo-500 to-purple-600 dark:from-blue-400 dark:via-blue-600 dark:to-indigo-500 drop-shadow-sm">
                Start Defending.
              </span>
            </h1>
            
            <p className="text-xl md:text-3xl text-zinc-600 dark:text-zinc-400 mb-16 leading-relaxed max-w-4xl mx-auto font-medium tracking-tight">
              Rampart is the unified network policy engine that thinks. Zero-trust orchestration for eBPF, nftables, and every Cloud provider in a single, autonomous control plane.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-8 mb-32">
              <button className="w-full sm:w-auto bg-blue-600 hover:bg-blue-500 text-white px-16 py-7 rounded-3xl font-black shadow-[0_30px_80px_-10px_rgba(37,99,235,0.4)] flex items-center justify-center gap-4 transition-all hover:-translate-y-2 active:scale-95 text-xl">
                Get Started Now
                <ChevronRight className="w-6 h-6" />
              </button>
              <button className="w-full sm:w-auto glass text-zinc-900 dark:text-white px-16 py-7 rounded-3xl font-black flex items-center justify-center gap-4 transition-all hover:bg-white dark:hover:bg-white/10 text-xl group border-2">
                <Github className="w-6 h-6 text-zinc-500 group-hover:text-blue-500 transition-colors" />
                View documentation
              </button>
            </div>
          </motion.div>
        </div>

        {/* Cinematic Dashboard Preview */}
        <motion.div 
          initial={{ opacity: 0, y: 100, scale: 0.9 }} animate={{ opacity: 1, y: 0, scale: 1 }} transition={{ delay: 0.5, duration: 1.2, ease: "circOut" }}
          className="relative max-w-6xl mx-auto"
        >
          <div className="glass rounded-[4rem] overflow-hidden shadow-[0_50px_150px_-20px_rgba(0,0,0,0.5)] border-zinc-200 dark:border-white/10 bg-black/40 backdrop-blur-3xl">
             <div className="bg-zinc-100/50 dark:bg-zinc-900/50 px-10 py-6 border-b border-zinc-200 dark:border-white/10 flex items-center justify-between">
                <div className="flex gap-3">
                  <div className="w-4 h-4 rounded-full bg-rose-500/20 border border-rose-500/40" />
                  <div className="w-4 h-4 rounded-full bg-amber-500/20 border border-amber-500/40" />
                  <div className="w-4 h-4 rounded-full bg-emerald-500/20 border border-emerald-500/40" />
                </div>
                <div className="flex items-center gap-3 text-zinc-500">
                   <Terminal className="w-4 h-4" />
                   <span className="text-[11px] font-black uppercase tracking-[0.5em] font-mono">policy-orchestrator-v0.1.0</span>
                </div>
                <div className="w-12" />
             </div>
             <div className="p-12 md:p-20 font-mono text-lg md:text-xl text-left">
                <pre className="text-zinc-800 dark:text-zinc-300 leading-[1.6] overflow-x-auto whitespace-pre">
                  <code className="block">
                    <span className="text-blue-500">apiVersion:</span> <span className="text-emerald-500">rampartfw.com/v1</span>{`\n`}
                    <span className="text-blue-500">kind:</span> <span className="text-emerald-500">PolicySet</span>{`\n`}
                    <span className="text-blue-500">metadata:</span>{`\n`}
                    {`  `} <span className="text-blue-500">name:</span> <span className="text-amber-500">"global-sentinel"</span>{`\n\n`}
                    <span className="text-zinc-500"># Autonomous Security Logic</span>{`\n`}
                    <span className="text-blue-500">policies:</span>{`\n`}
                    {`  - `}<span className="text-blue-500">name:</span> <span className="text-amber-500">"zero-trust-edge"</span>{`\n`}
                    {`    `}<span className="text-blue-500">rules:</span>{`\n`}
                    {`      - `}<span className="text-blue-500">name:</span> <span className="text-amber-500">"mitigate-l7-threats"</span>{`\n`}
                    {`        `}<span className="text-blue-500">action:</span> <span className="text-rose-500 font-black">drop</span>{`\n`}
                    {`        `}<span className="text-blue-500">match:</span>{`\n`}
                    {`          `}<span className="text-blue-500">dns:</span> <span className="text-zinc-400">{`{ query: "evil-botnet.io" }`}</span>
                  </code>
                </pre>
             </div>
          </div>
          {/* Floating UI Elements */}
          <motion.div 
            animate={{ y: [0, -20, 0] }} transition={{ duration: 4, repeat: Infinity }}
            className="absolute -right-12 top-1/4 glass p-6 rounded-3xl shadow-2xl border-white/20 hidden lg:block"
          >
             <Activity className="w-8 h-8 text-emerald-500 mb-2" />
             <div className="text-[10px] font-black uppercase tracking-widest text-zinc-500">Throughput</div>
             <div className="text-2xl font-black dark:text-white">1.2M <span className="text-xs opacity-50">pps</span></div>
          </motion.div>
          <motion.div 
            animate={{ y: [0, 20, 0] }} transition={{ duration: 5, repeat: Infinity, delay: 1 }}
            className="absolute -left-12 bottom-1/4 glass p-6 rounded-3xl shadow-2xl border-white/20 hidden lg:block"
          >
             <Lock className="w-8 h-8 text-blue-500 mb-2" />
             <div className="text-[10px] font-black uppercase tracking-widest text-zinc-500">mTLS Status</div>
             <div className="text-sm font-black dark:text-white uppercase tracking-tighter">🔒 Fully Synchronized</div>
          </motion.div>
        </motion.div>
      </div>
    </section>
  )
}

const FeatureGrid = () => {
  const features = [
    { icon: Cpu, title: "eBPF/XDP Fast-Path", desc: "Native kernel execution with microsecond latency. Offload massive rule sets to the network driver for unmatched performance.", color: "blue" },
    { icon: Bot, title: "Autonomous Sentinel", desc: "Beyond traditional IPS. Real-time threat scoring and otonom cluster-wide mitigation powered by DPI signals.", color: "indigo" },
    { icon: Globe, title: "Unified Cloud Plane", desc: "Single control plane for AWS, GCP, and Azure. Synchronize security groups globally with one YAML manifest.", color: "purple" },
    { icon: RefreshCw, title: "Self-Healing Watchdog", desc: "Continuous drift detection. Rampart automatically audits and corrects firewall state to ensure 100% compliance.", color: "emerald" },
    { icon: Fingerprint, title: "Mühürlü Audit Chain", desc: "Cryptographic hash-chaining for all logs. Tamper-evident, high-integrity audit trails for mission-critical compliance.", color: "rose" },
    { icon: Search, title: "Layer-7 DPI Intelligence", desc: "Deep Packet Inspection for DNS, HTTP, and TLS SNI. Stop application-level threats before they enter your network.", color: "amber" },
  ]

  return (
    <section id="features" className="py-48 bg-white dark:bg-zinc-950 relative overflow-hidden">
      <div className="container mx-auto px-8">
        <div className="flex flex-col md:flex-row items-end justify-between gap-8 mb-24">
           <div className="max-w-2xl text-left">
              <h2 className="text-5xl md:text-7xl font-black mb-8 tracking-tighter dark:text-white text-zinc-950 leading-none">
                 The Edge of <br />
                 <span className="text-blue-600 underline decoration-8 underline-offset-8">Infrastructure.</span>
              </h2>
              <p className="text-xl text-zinc-500 dark:text-zinc-400 font-bold leading-relaxed">
                 Traditional firewalls are static. Rampart is dynamic, distributed, and aware. Built for the modern cloud-native era.
              </p>
           </div>
           <div className="hidden lg:block pb-4">
              <button className="px-8 py-4 rounded-2xl glass font-black text-sm uppercase tracking-widest hover:border-blue-500/50 transition-all flex items-center gap-3">
                 View benchmark suite <ArrowRight className="w-4 h-4" />
              </button>
           </div>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-10">
           {features.map((f, i) => (
             <motion.div 
               key={i} whileHover={{ y: -16 }}
               className="p-12 rounded-[3.5rem] premium-card relative overflow-hidden group"
             >
                <div className={cn(
                  "w-16 h-16 rounded-2xl flex items-center justify-center mb-10 transition-all duration-500 scale-110 shadow-2xl",
                  `bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-400 group-hover:bg-blue-600 group-hover:text-white group-hover:shadow-blue-600/30`
                )}>
                   <f.icon className="w-8 h-8" />
                </div>
                <h3 className="text-3xl font-black mb-6 dark:text-white text-zinc-950 tracking-tighter leading-none">{f.title}</h3>
                <p className="text-zinc-500 dark:text-zinc-400 text-lg leading-relaxed font-bold opacity-80">{f.desc}</p>
                <div className="absolute -bottom-8 -right-8 w-32 h-32 bg-blue-600/5 blur-3xl rounded-full opacity-0 group-hover:opacity-100 transition-opacity" />
             </motion.div>
           ))}
        </div>
      </div>
    </section>
  )
}

const DocumentationPortal = () => {
  const [activeTab, setActiveTab] = useState('install')
  
  const docs = [
    {
      id: 'install',
      title: "Quick Deployment",
      icon: TerminalSquare,
      content: (
        <div className="space-y-10 animate-in fade-in slide-in-from-right-4 duration-500">
          <div className="space-y-4">
             <h4 className="text-3xl font-black dark:text-white text-zinc-950 tracking-tighter underline decoration-blue-600 decoration-4 underline-offset-4">One-Command Setup.</h4>
             <p className="text-lg text-zinc-600 dark:text-zinc-400 font-medium">Deploy the global controller and local node with our automated installer.</p>
          </div>
          <div className="bg-zinc-950 p-10 rounded-[2.5rem] border border-white/10 font-mono shadow-inner group relative overflow-hidden">
             <div className="flex items-center gap-2 text-zinc-600 mb-4 text-xs font-bold uppercase tracking-widest">
                <Globe className="w-3 h-3 text-emerald-500" />
                Secure Download Channel
             </div>
             <div className="text-xl md:text-2xl font-black flex items-center gap-4">
                <span className="text-zinc-500 select-none">$ </span>
                <span className="text-emerald-400">curl -sSL rampartfw.com/install | sh</span>
             </div>
             <div className="absolute top-0 right-0 w-64 h-full bg-gradient-to-l from-emerald-500/5 to-transparent pointer-events-none" />
          </div>
          <div className="grid md:grid-cols-2 gap-6">
             <div className="p-8 rounded-3xl bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10">
                <CheckCircle2 className="w-6 h-6 text-blue-500 mb-4" />
                <h5 className="font-black text-lg mb-2">Automated Discovery</h5>
                <p className="text-sm text-zinc-500 font-medium">Detects OS, Arch, and kernel features (XDP/BTF) automatically.</p>
             </div>
             <div className="p-8 rounded-3xl bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10">
                <CheckCircle2 className="w-6 h-6 text-emerald-500 mb-4" />
                <h5 className="font-black text-lg mb-2">Systemd Ready</h5>
                <p className="text-sm text-zinc-500 font-medium">Hooks into init systems with pre-configured self-healing services.</p>
             </div>
          </div>
        </div>
      )
    },
    {
      id: 'architecture',
      title: "Control Plane Architecture",
      icon: Network,
      content: (
        <div className="space-y-10 animate-in fade-in slide-in-from-right-4 duration-500">
           <div className="space-y-4">
             <h4 className="text-3xl font-black dark:text-white text-zinc-950 tracking-tighter underline decoration-indigo-600 decoration-4 underline-offset-4">Distributed Intelligence.</h4>
             <p className="text-lg text-zinc-600 dark:text-zinc-400 font-medium">Rampart is designed for zero-single-point-of-failure clusters.</p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
             <div className="space-y-4">
                <div className="p-4 rounded-xl bg-blue-600/10 text-blue-600 font-black text-[10px] uppercase tracking-widest inline-block">The Fast-Path</div>
                <p className="text-zinc-500 font-medium leading-relaxed">Leverages eBPF/XDP tail calls to process traffic before it even hits the Linux networking stack. Perfect for high-volume DDoS mitigation.</p>
             </div>
             <div className="space-y-4">
                <div className="p-4 rounded-xl bg-purple-600/10 text-purple-600 font-black text-[10px] uppercase tracking-widest inline-block">The Multi-Cloud Plane</div>
                <p className="text-zinc-500 font-medium leading-relaxed">Translates abstract intent into AWS Security Groups, GCP Firewall Rules, and Azure NSGs simultaneously.</p>
             </div>
          </div>
          <div className="p-8 rounded-3xl glass border-indigo-500/20 text-center relative overflow-hidden">
             <Server className="w-12 h-12 text-indigo-500 mx-auto mb-4" />
             <h5 className="font-black text-xl mb-2 italic">Raft Consensus Core</h5>
             <p className="text-sm text-zinc-500 font-bold tracking-tight px-10">Ensures all nodes converge to the exact same policy state using mTLS mühürlü communication.</p>
             <div className="absolute -z-10 top-0 left-0 w-full h-full bg-indigo-500/5 blur-3xl rounded-full" />
          </div>
        </div>
      )
    },
    {
      id: 'sentinel',
      title: "Autonomous sentinel",
      icon: Bot,
      content: (
        <div className="space-y-8 animate-in fade-in slide-in-from-right-4 duration-500">
           <div className="space-y-4">
             <h4 className="text-3xl font-black dark:text-white text-zinc-950 tracking-tighter underline decoration-rose-600 decoration-4 underline-offset-4">Otonom Threat Response.</h4>
             <p className="text-lg text-zinc-600 dark:text-zinc-400 font-medium">Stop threats in milliseconds without human intervention.</p>
          </div>
          <div className="p-10 rounded-[3rem] bg-zinc-950 border border-white/10 shadow-2xl relative overflow-hidden">
             <div className="flex items-center gap-4 mb-8">
                <div className="w-3 h-3 rounded-full bg-rose-500 animate-pulse" />
                <span className="text-[10px] font-black uppercase tracking-[0.4em] text-rose-500">Sentinel: Risk Score Reached Threshold</span>
             </div>
             <div className="space-y-4 font-mono text-sm mb-10">
                <p className="text-zinc-500"># Real-time traffic analysis result</p>
                <div className="grid grid-cols-2 gap-4">
                   <div className="p-4 rounded-xl bg-white/5 border border-white/10">
                      <p className="text-zinc-600 text-[10px] mb-1 uppercase font-bold">Source IP</p>
                      <p className="text-emerald-400">192.168.1.45</p>
                   </div>
                   <div className="p-4 rounded-xl bg-white/5 border border-white/10">
                      <p className="text-zinc-600 text-[10px] mb-1 uppercase font-bold">Signal</p>
                      <p className="text-rose-400 font-black italic">DNS_FLOOD</p>
                   </div>
                </div>
                <p className="text-zinc-300 italic">{`> SCORING: IP 192.168.1.45 -> RISK 85 (Target 70)`}</p>
                <p className="text-emerald-400 font-black underline">{`> REACTION: Global Block Applied via Raft.`}</p>
             </div>
             <button className="text-blue-500 text-xs font-black uppercase tracking-[0.5em] flex items-center gap-3">
                Configure Sentinel Intelligence <ArrowRight className="w-4 h-4" />
             </button>
             <div className="absolute top-0 right-0 w-64 h-full bg-rose-600/5 blur-[100px] pointer-events-none" />
          </div>
        </div>
      )
    },
    {
        id: 'ai',
        title: "AI-Ready (MCP)",
        icon: Command,
        content: (
          <div className="space-y-10 animate-in fade-in slide-in-from-right-4 duration-500">
             <div className="space-y-4">
               <h4 className="text-3xl font-black dark:text-white text-zinc-950 tracking-tighter underline decoration-emerald-600 decoration-4 underline-offset-4">Agentic Orchestration.</h4>
               <p className="text-lg text-zinc-600 dark:text-zinc-400 font-medium">The first firewall with native Model Context Protocol support.</p>
            </div>
            <div className="grid md:grid-cols-2 gap-10">
                <div className="p-8 rounded-3xl glass border-emerald-500/20">
                    <Bot className="w-10 h-10 text-emerald-500 mb-6" />
                    <h5 className="text-xl font-black mb-3">Claude & GPT Ready</h5>
                    <p className="text-sm text-zinc-500 font-bold leading-relaxed">Connect your favorite AI agent directly to Rampart. Let it analyze traffic, suggest policies, and respond to incidents.</p>
                </div>
                <div className="p-8 rounded-3xl glass border-blue-500/20">
                    <TerminalSquare className="w-10 h-10 text-blue-500 mb-6" />
                    <h5 className="text-xl font-black mb-3">Tool Call Interface</h5>
                    <p className="text-sm text-zinc-500 font-bold leading-relaxed">Exposes high-level tools like <code className="text-blue-500">apply_policy</code>, <code className="text-blue-500">list_rules</code>, and <code className="text-blue-500">simulate_packet</code> to LLMs.</p>
                </div>
            </div>
          </div>
        )
      }
  ]

  return (
    <section id="docs" className="py-48 bg-zinc-50 dark:bg-zinc-900/10 border-y border-zinc-200 dark:border-white/5 relative overflow-hidden">
      <div className="container mx-auto px-8 relative z-10">
        <div className="max-w-7xl mx-auto">
          <div className="flex flex-col lg:flex-row gap-24">
            <div className="lg:w-[350px] shrink-0">
              <h2 className="text-6xl font-black mb-10 tracking-tighter dark:text-white text-zinc-950 underline decoration-blue-600 decoration-[16px] underline-offset-8">Docs.</h2>
              <p className="text-zinc-500 dark:text-zinc-400 text-2xl mb-16 leading-tight font-black tracking-tighter">
                Master the fortress. Orchestrate the edge.
              </p>
              
              <div className="space-y-4">
                {docs.map((doc) => (
                  <button 
                    key={doc.id}
                    onClick={() => setActiveId(doc.id)}
                    className={cn(
                      "w-full flex items-center justify-between px-10 py-6 rounded-[2rem] font-black text-[12px] uppercase tracking-[0.4em] transition-all",
                      activeId === doc.id 
                        ? "bg-blue-600 text-white shadow-2xl shadow-blue-600/40 translate-x-6 scale-105" 
                        : "text-zinc-400 dark:text-zinc-600 hover:bg-zinc-100 dark:hover:bg-white/5"
                    )}
                  >
                    <div className="flex items-center gap-5">
                       <doc.icon className="w-5 h-5" />
                       {doc.title}
                    </div>
                    <ChevronRight className={cn("w-4 h-4 transition-transform", activeId === doc.id ? "rotate-90" : "")} />
                  </button>
                ))}
              </div>
            </div>

            <div className="flex-1 glass rounded-[5rem] p-12 md:p-24 border-zinc-200 dark:border-white/10 bg-white dark:bg-black/40 relative shadow-[0_60px_100px_-20px_rgba(0,0,0,0.4)]">
               <AnimatePresence mode="wait">
                 <motion.div
                   key={activeId}
                   initial={{ opacity: 0, x: 40 }} animate={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: -40 }}
                   transition={{ duration: 0.5, ease: "circOut" }}
                   className="relative z-10"
                 >
                   {docs.find(d => d.id === activeId)?.content}
                 </motion.div>
               </AnimatePresence>
               <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-blue-600/5 blur-[120px] rounded-full pointer-events-none" />
            </div>
          </div>
        </div>
      </div>
      <div className="absolute bottom-0 left-0 w-full h-64 bg-gradient-to-t from-white dark:from-zinc-950 to-transparent pointer-events-none" />
    </section>
  )
}

const Releases = () => {
  const versions = [
    { v: 'v0.1.0', tag: 'Initial Beta (Orchestrator)', date: 'April 2026', current: true, features: ['Multi-Backend', 'Raft Cluster', 'IPS Core', 'eBPF/XDP'] },
    { v: 'v0.1.1', tag: 'Maintenance', date: 'May 2026', current: false, features: ['WebUI Polish', 'GCP/Azure fixes'] },
  ]

  return (
    <section id="releases" className="py-48 bg-white dark:bg-zinc-950">
      <div className="container mx-auto px-8">
        <div className="max-w-5xl mx-auto">
           <div className="flex items-center gap-6 mb-20">
              <div className="w-16 h-16 bg-zinc-100 dark:bg-white/5 rounded-3xl flex items-center justify-center border border-zinc-200 dark:border-white/10 shadow-inner">
                 <Box className="w-8 h-8 text-blue-600" />
              </div>
              <h2 className="text-6xl md:text-7xl font-black tracking-tighter dark:text-white text-zinc-950">Releases.</h2>
           </div>
           
           <div className="space-y-12">
              {versions.map(version => (
                <div key={version.v} className="p-12 md:p-16 rounded-[4rem] glass border-zinc-200 dark:border-white/10 flex flex-col md:flex-row md:items-start justify-between gap-12 group transition-all duration-500 hover:border-blue-500/30">
                   <div className="flex gap-10">
                      <div className={cn(
                        "w-24 h-24 rounded-[2rem] flex items-center justify-center font-black text-3xl shadow-2xl shrink-0 transition-transform duration-700 group-hover:rotate-12",
                        version.current ? "bg-blue-600 text-white shadow-blue-600/30" : "bg-zinc-100 dark:bg-zinc-800 text-zinc-400"
                      )}>
                        {version.v.split('.')[1]}.{version.v.split('.')[2]}
                      </div>
                      <div className="space-y-4">
                         <h4 className="text-4xl font-black dark:text-white text-zinc-950 flex items-center gap-5 tracking-tighter leading-none">
                           {version.v}
                           {version.current && <span className="px-5 py-2 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 text-[10px] font-black uppercase tracking-[0.3em]">Stable Release</span>}
                         </h4>
                         <p className="text-xl text-zinc-500 font-bold tracking-tight leading-none mb-6">{version.tag} • {version.date}</p>
                         <div className="flex flex-wrap gap-3 pt-4">
                            {version.features.map(f => (
                               <span key={f} className="px-4 py-2 rounded-xl bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 text-xs font-black dark:text-zinc-400 uppercase tracking-widest">{f}</span>
                            ))}
                         </div>
                      </div>
                   </div>
                   
                   <div className="flex gap-4 shrink-0">
                      <a 
                        href={`https://github.com/ersinkoc/Rampart/releases/tag/${version.v}`}
                        className="px-12 py-5 rounded-2xl bg-zinc-950 dark:bg-white text-white dark:text-zinc-950 font-black text-sm flex items-center gap-3 hover:scale-110 active:scale-95 transition-all shadow-xl shadow-black/20"
                      >
                         <Download className="w-5 h-5" />
                         Get Assets
                      </a>
                   </div>
                </div>
              ))}
           </div>
        </div>
      </div>
    </section>
  )
}

const App = () => {
  return (
    <ThemeProvider>
      <main className="bg-white dark:bg-zinc-950 min-h-screen text-zinc-900 dark:text-white selection:bg-blue-500/30 font-sans antialiased overflow-x-hidden">
        <Navbar />
        <Hero />
        <FeatureGrid />
        <DocumentationPortal />
        <Releases />
        
        <section id="architecture" className="py-48 container mx-auto px-8 text-center relative">
          <div className="glass p-20 md:p-48 rounded-[6rem] border-zinc-200 dark:border-white/10 relative overflow-hidden bg-zinc-50 dark:bg-zinc-900/60 shadow-[0_80px_200px_-20px_rgba(0,0,0,0.5)]">
            <h2 className="text-7xl md:text-[12rem] font-black mb-16 relative z-10 tracking-[-0.07em] leading-[0.75] dark:text-white text-zinc-950">
              Future <br />
              <span className="text-blue-600 drop-shadow-[0_20px_50px_rgba(37,99,235,0.3)]">Mühürlendi.</span>
            </h2>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-10 relative z-10">
              <a href="https://github.com/ersinkoc/Rampart" className="w-full sm:w-auto bg-zinc-950 dark:bg-white text-white dark:text-zinc-950 px-24 py-10 rounded-[2.5rem] font-black hover:scale-110 transition-all text-3xl shadow-[0_30px_100px_-10px_rgba(0,0,0,0.8)] dark:shadow-white/10 active:scale-95">
                Fortify Now v0.1.0
              </a>
            </div>
            <div className="absolute top-0 right-0 w-[90%] h-full bg-[radial-gradient(circle_at_top_right,rgba(59,130,246,0.2),transparent)]" />
            <div className="absolute bottom-0 left-0 w-[60%] h-full bg-[radial-gradient(circle_at_bottom_left,rgba(99,102,241,0.1),transparent)]" />
          </div>
        </section>

        <footer className="py-32 border-t border-zinc-200 dark:border-white/5 bg-zinc-50 dark:bg-black/90">
          <div className="container mx-auto px-8 flex flex-col md:flex-row items-center justify-between gap-20">
            <div className="flex items-center gap-6 group cursor-pointer shrink-0">
              <div className="w-16 h-16 bg-blue-600 rounded-[1.5rem] flex items-center justify-center shadow-3xl shadow-blue-600/40 group-hover:rotate-[-12deg] transition-all duration-700">
                <Shield className="text-white w-9 h-9" />
              </div>
              <div className="flex flex-col">
                <span className="text-3xl font-black tracking-tighter dark:text-white text-zinc-950 leading-none">RAMPART</span>
                <span className="text-[10px] font-black uppercase tracking-[0.6em] text-blue-600 mt-2">TACTICAL INTELLIGENCE</span>
              </div>
            </div>
            
            <p className="text-zinc-500 font-bold text-lg max-w-sm text-center md:text-left leading-relaxed opacity-70">
              High-performance network orchestration. Secure by design. Autonomous by default.
            </p>
            
            <div className="flex items-center gap-12 shrink-0">
              <a href="https://github.com/ersinkoc/Rampart" className="text-zinc-400 hover:text-blue-600 transition-all scale-[2] hover:rotate-12">
                <Github />
              </a>
              <a href="#" className="text-zinc-400 hover:text-blue-600 transition-all scale-[2] hover:rotate-[-12deg]">
                <Globe />
              </a>
            </div>
          </div>
          <div className="container mx-auto px-8 mt-24 pt-12 border-t border-zinc-200 dark:border-white/5 text-center">
             <p className="text-zinc-600 dark:text-zinc-500 text-[10px] font-black uppercase tracking-[0.5em]">
                Licensed under Apache 2.0 • Engineered for the Global Edge
             </p>
          </div>
        </footer>
      </main>
    </ThemeProvider>
  )
}

export default App
