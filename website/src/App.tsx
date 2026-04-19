import { useState, useEffect, createContext, useContext } from 'react'
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
  ArrowRight
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
  const [isDark, setIsDark] = useState(true)

  useEffect(() => {
    const root = window.document.documentElement
    if (isDark) {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
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
    { name: 'Docs', href: '#docs' },
    { name: 'Architecture', href: '#architecture' },
  ]

  return (
    <nav className={cn(
      "fixed top-0 left-0 right-0 z-50 transition-all duration-300 border-b",
      isScrolled ? "glass py-3 border-zinc-200 dark:border-zinc-800/50" : "bg-transparent py-5 border-transparent"
    )}>
      <div className="container mx-auto px-6 flex items-center justify-between">
        <div className="flex items-center gap-2 group cursor-pointer" onClick={() => window.scrollTo({top: 0, behavior: 'smooth'})}>
          <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center group-hover:rotate-12 transition-transform shadow-lg shadow-blue-500/20">
            <Shield className="text-white w-6 h-6" />
          </div>
          <span className="text-xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-zinc-900 via-zinc-700 to-zinc-500 dark:from-white dark:to-zinc-400">
            RAMPART
          </span>
        </div>

        <div className="hidden md:flex items-center gap-8">
          {navLinks.map(link => (
            <a key={link.name} href={link.href} className="text-sm font-medium text-zinc-500 dark:text-zinc-400 hover:text-blue-600 dark:hover:text-white transition-colors">
              {link.name}
            </a>
          ))}
          <div className="h-4 w-[1px] bg-zinc-200 dark:bg-zinc-800" />
          <button 
            onClick={toggle}
            className="p-2 rounded-lg bg-zinc-100 dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400 hover:text-blue-600 dark:hover:text-white transition-all"
          >
            {isDark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
          </button>
          <a 
            href="https://github.com/ersinkoc/Rampart" 
            className="flex items-center gap-2 bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 hover:opacity-90 px-4 py-2 rounded-lg text-sm font-bold transition-all"
          >
            <Github className="w-4 h-4" />
            GitHub
          </a>
        </div>

        <button className="md:hidden text-zinc-500 dark:text-zinc-400" onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
          {isMobileMenuOpen ? <X /> : <Menu />}
        </button>
      </div>

      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div 
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="md:hidden absolute top-full left-0 right-0 glass border-b border-zinc-200 dark:border-zinc-800 flex flex-col p-6 gap-4"
          >
            {navLinks.map(link => (
              <a key={link.name} href={link.href} className="text-lg font-medium dark:text-zinc-300" onClick={() => setIsMobileMenuOpen(false)}>
                {link.name}
              </a>
            ))}
            <button onClick={() => { toggle(); setIsMobileMenuOpen(false); }} className="flex items-center gap-2 text-lg font-medium dark:text-zinc-300">
              {isDark ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
              {isDark ? 'Light Mode' : 'Dark Mode'}
            </button>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  )
}

const Hero = () => {
  return (
    <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full hero-gradient pointer-events-none" />
      <div className="container mx-auto px-6 relative z-10 text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            <span className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-600 dark:text-blue-400 text-xs font-bold uppercase tracking-widest mb-6">
              <Zap className="w-3 h-3 fill-current" />
              v0.1.0 Initial Beta Release
            </span>
            <h1 className="text-5xl md:text-8xl font-extrabold mb-8 tracking-tight leading-tight dark:text-white">
              The Fortress for <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-600 via-blue-500 to-indigo-400">
                Modern Networks
              </span>
            </h1>
            <p className="text-lg md:text-xl text-zinc-600 dark:text-zinc-400 mb-10 leading-relaxed max-w-2xl mx-auto">
              Abstract the complexity of Linux eBPF/XDP and Cloud security groups behind a single, intelligent YAML interface. Deploy anywhere, secure everything.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <button className="w-full sm:w-auto bg-blue-600 hover:bg-blue-500 text-white px-10 py-4 rounded-2xl font-bold shadow-xl shadow-blue-500/25 flex items-center justify-center gap-2 group transition-all">
                Get Started Now
                <ChevronRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </button>
              <button className="w-full sm:w-auto bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-800 text-zinc-900 dark:text-white px-10 py-4 rounded-2xl font-bold flex items-center justify-center gap-2 transition-all">
                <Github className="w-5 h-5" />
                Star on GitHub
              </button>
            </div>
          </motion.div>

        <motion.div 
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.3, duration: 0.8 }}
          className="mt-20 max-w-5xl mx-auto glass rounded-3xl overflow-hidden shadow-2xl border-zinc-200 dark:border-zinc-800"
        >
          <div className="bg-zinc-100 dark:bg-zinc-900/50 px-6 py-4 border-b border-zinc-200 dark:border-zinc-800 flex items-center justify-between">
            <div className="flex gap-2">
              <div className="w-3 h-3 rounded-full bg-rose-500/20 border border-rose-500/40" />
              <div className="w-3 h-3 rounded-full bg-amber-500/20 border border-amber-500/40" />
              <div className="w-3 h-3 rounded-full bg-emerald-500/20 border border-emerald-500/40" />
            </div>
            <span className="text-[10px] text-zinc-500 font-mono tracking-widest uppercase flex items-center gap-2">
              <Code2 className="w-3 h-3" />
              main-ingress.yaml
            </span>
            <div className="w-10" />
          </div>
          <div className="p-8 md:p-12 font-mono text-sm md:text-base text-left bg-white dark:bg-zinc-950/50">
            <pre className="text-zinc-800 dark:text-zinc-300 leading-relaxed overflow-x-auto">
              <code>{`apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: "global-ingress"

policies:
  - name: "autonomous-protection"
    priority: 100
    rules:
      - name: "block-malicious-dns"
        action: drop
        match:
          appProtocol: dns
          dns: { query: "malicious.com" }
      
      - name: "high-perf-web-access"
        action: accept
        match:
          protocol: tcp
          destPorts: [80, 443]`}</code>
            </pre>
          </div>
        </motion.div>
      </div>
    </section>
  )
}

const Docs = () => {
  const sections = [
    {
      title: "Getting Started",
      icon: Book,
      content: [
        { name: "Installation", desc: "Single binary deployment for Linux and macOS." },
        { name: "First Policy", desc: "Writing your first YAML-based network rule." },
        { name: "Apply & Plan", desc: "Deterministic deployments with dry-run support." }
      ]
    },
    {
      title: "Core Concepts",
      icon: Layers,
      content: [
        { name: "Policy Engine", desc: "How Rampart compiles abstract YAML to kernel instructions." },
        { name: "Distributed Raft", desc: "Ensuring 100% consistency across your entire fleet." },
        { name: "Autonomous IPS", desc: "Smart threat detection with dynamic risk scoring." }
      ]
    },
    {
      title: "Advanced",
      icon: Cpu,
      content: [
        { name: "eBPF/XDP", desc: "Hyper-performance filtering at the NIC driver level." },
        { name: "Cloud Backends", desc: "Orchestrating AWS, GCP, and Azure from one place." },
        { name: "MCP Server", desc: "Automating your firewall with AI-driven agents." }
      ]
    }
  ]

  return (
    <section id="docs" className="py-24 bg-white dark:bg-zinc-950">
      <div className="container mx-auto px-6">
        <div className="flex flex-col md:flex-row gap-16">
          <div className="md:w-1/3">
            <h2 className="text-4xl font-bold mb-6 tracking-tight dark:text-white">Documentation</h2>
            <p className="text-zinc-600 dark:text-zinc-400 text-lg mb-8 leading-relaxed">
              Explore our comprehensive guides to mastering Rampart and securing your infrastructure.
            </p>
            <div className="p-6 rounded-2xl bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800">
               <div className="flex items-center gap-3 mb-4">
                 <Terminal className="w-5 h-5 text-blue-600" />
                 <span className="font-bold dark:text-white">Quick Install</span>
               </div>
               <code className="text-sm font-mono block bg-zinc-200 dark:bg-zinc-950 p-3 rounded-lg dark:text-zinc-300">
                 curl -sSL rampartfw.com/install | sh
               </code>
            </div>
          </div>
          <div className="md:w-2/3 grid sm:grid-cols-2 gap-10">
            {sections.map((section, idx) => (
              <div key={idx} className="space-y-6">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-blue-600/10 flex items-center justify-center text-blue-600">
                    <section.icon className="w-5 h-5" />
                  </div>
                  <h3 className="text-xl font-bold dark:text-white">{section.title}</h3>
                </div>
                <ul className="space-y-4">
                  {section.content.map((item, i) => (
                    <li key={i} className="group cursor-pointer">
                      <h4 className="font-bold text-zinc-900 dark:text-zinc-200 group-hover:text-blue-600 transition-colors flex items-center gap-2">
                        {item.name}
                        <ArrowRight className="w-4 h-4 opacity-0 group-hover:opacity-100 -translate-x-2 group-hover:translate-x-0 transition-all" />
                      </h4>
                      <p className="text-sm text-zinc-500 dark:text-zinc-400 mt-1">{item.desc}</p>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}

const FeatureCard = ({ icon: Icon, title, description, delay }: { icon: any, title: string, description: string, delay: number }) => (
  <motion.div 
    initial={{ opacity: 0, y: 20 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true }}
    transition={{ delay, duration: 0.5 }}
    className="p-8 rounded-3xl bg-zinc-50 dark:bg-zinc-900/30 border border-zinc-200 dark:border-zinc-800 hover:border-blue-500/30 hover:bg-blue-500/[0.02] transition-all group"
  >
    <div className="w-12 h-12 bg-white dark:bg-zinc-800 rounded-2xl flex items-center justify-center mb-6 shadow-sm group-hover:bg-blue-600 group-hover:text-white transition-all">
      <Icon className="w-6 h-6 text-blue-600 dark:text-blue-400 group-hover:text-white transition-colors" />
    </div>
    <h3 className="text-xl font-bold mb-3 text-zinc-900 dark:text-white">{title}</h3>
    <p className="text-zinc-600 dark:text-zinc-400 leading-relaxed">{description}</p>
  </motion.div>
)

const App = () => {
  return (
    <ThemeProvider>
      <main className="bg-white dark:bg-zinc-950 min-h-screen text-zinc-900 dark:text-white selection:bg-blue-500/30">
        <Navbar />
        <Hero />
        
        <section id="features" className="py-24 border-y border-zinc-200 dark:border-zinc-900">
          <div className="container mx-auto px-6">
            <div className="text-center max-w-2xl mx-auto mb-20">
              <h2 className="text-3xl md:text-5xl font-bold mb-6 tracking-tight">Built for Resilience.</h2>
              <p className="text-zinc-600 dark:text-zinc-400 text-lg leading-relaxed">
                Rampart is more than a firewall manager. It's an autonomous sentinel that defends your fleet 24/7.
              </p>
            </div>
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
              <FeatureCard icon={Cpu} title="eBPF/XDP" description="Microsecond latency filtering at the driver level." delay={0.1} />
              <FeatureCard icon={Bot} title="Sentinel IPS" description="Otonom threat detection and cluster-wide response." delay={0.2} />
              <FeatureCard icon={Cloud} title="Unified Cloud" description="Manage AWS, GCP, and Azure from a single YAML." delay={0.3} />
              <FeatureCard icon={RefreshCw} title="Self-Healing" description="Continuous watchdog monitoring for zero drift." delay={0.4} />
              <FeatureCard icon={Search} title="DPI Engine" description="Layer-7 awareness for DNS and HTTP traffic." delay={0.5} />
              <FeatureCard icon={Globe} title="mTLS Raft" description="Secure, distributed consensus for every node." delay={0.6} />
            </div>
          </div>
        </section>

        <Docs />

        <section id="architecture" className="py-24 container mx-auto px-6 text-center">
          <div className="glass p-12 md:p-24 rounded-[3.5rem] border-zinc-200 dark:border-zinc-800 relative overflow-hidden bg-zinc-50 dark:bg-zinc-900/50">
            <h2 className="text-4xl md:text-6xl font-extrabold mb-8 relative z-10 tracking-tight">Ready to fortify your <br /><span className="text-blue-600">infrastructure?</span></h2>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 relative z-10">
              <a href="https://github.com/ersinkoc/Rampart" className="w-full sm:w-auto bg-blue-600 text-white px-12 py-5 rounded-2xl font-bold hover:bg-blue-500 transition-all text-lg shadow-2xl shadow-blue-500/20">
                Get Started v0.1.0
              </a>
              <a href="#docs" className="w-full sm:w-auto bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 px-12 py-5 rounded-2xl font-bold hover:opacity-90 transition-all text-lg">
                Read the Docs
              </a>
            </div>
          </div>
        </section>

        <footer className="py-12 border-t border-zinc-200 dark:border-zinc-900">
          <div className="container mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-8">
            <div className="flex items-center gap-2">
              <Shield className="text-blue-600 w-6 h-6" />
              <span className="text-lg font-bold tracking-tight">RAMPART</span>
            </div>
            <p className="text-zinc-500 text-sm">
              © 2026 Rampart Intelligence. Licensed under Apache 2.0.
            </p>
            <div className="flex items-center gap-6">
              <a href="https://github.com/ersinkoc/Rampart" className="text-zinc-400 hover:text-blue-600 transition-colors">
                <Github className="w-5 h-5" />
              </a>
              <a href="#" className="text-zinc-400 hover:text-blue-600 transition-colors">
                <Globe className="w-5 h-5" />
              </a>
            </div>
          </div>
        </footer>
      </main>
    </ThemeProvider>
  )
}

export default App
