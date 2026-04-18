import { useState, useEffect } from 'react'
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
  CheckCircle2
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const Navbar = () => {
  const [isScrolled, setIsScrolled] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20)
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  const navLinks = [
    { name: 'Features', href: '#features' },
    { name: 'Architecture', href: '#architecture' },
    { name: 'Documentation', href: 'https://github.com/ersinkoc/Rampart#documentation' },
  ]

  return (
    <nav className={cn(
      "fixed top-0 left-0 right-0 z-50 transition-all duration-300 border-b",
      isScrolled ? "glass py-3 border-zinc-800/50" : "bg-transparent py-5 border-transparent"
    )}>
      <div className="container mx-auto px-6 flex items-center justify-between">
        <div className="flex items-center gap-2 group cursor-pointer">
          <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center group-hover:rotate-12 transition-transform shadow-lg shadow-blue-500/20">
            <Shield className="text-white w-6 h-6" />
          </div>
          <span className="text-xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-zinc-400">
            RAMPART
          </span>
        </div>

        <div className="hidden md:flex items-center gap-8">
          {navLinks.map(link => (
            <a key={link.name} href={link.href} className="text-sm font-medium text-zinc-400 hover:text-white transition-colors">
              {link.name}
            </a>
          ))}
          <a 
            href="https://github.com/ersinkoc/Rampart" 
            className="flex items-center gap-2 bg-zinc-800 hover:bg-zinc-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-all"
          >
            <Github className="w-4 h-4" />
            GitHub
          </a>
        </div>

        <button className="md:hidden text-zinc-400" onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
          {isMobileMenuOpen ? <X /> : <Menu />}
        </button>
      </div>

      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div 
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="md:hidden absolute top-full left-0 right-0 glass border-b border-zinc-800 flex flex-col p-6 gap-4"
          >
            {navLinks.map(link => (
              <a key={link.name} href={link.href} className="text-lg font-medium text-zinc-300" onClick={() => setIsMobileMenuOpen(false)}>
                {link.name}
              </a>
            ))}
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  )
}

const Hero = () => {
  return (
    <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden hero-gradient">
      <div className="container mx-auto px-6 relative z-10">
        <div className="max-w-4xl mx-auto text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            <span className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400 text-xs font-bold uppercase tracking-widest mb-6">
              <Zap className="w-3 h-3 fill-current" />
              v1.0.0 Now Production Ready
            </span>
            <h1 className="text-5xl md:text-7xl font-extrabold mb-8 tracking-tight leading-tight">
              Autonomous Network <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-blue-600 to-indigo-500">
                Defense & Orchestration
              </span>
            </h1>
            <p className="text-lg md:text-xl text-zinc-400 mb-10 leading-relaxed max-w-2xl mx-auto">
              Abstract the complexity of Linux eBPF/XDP and Cloud security groups behind a single, intelligent YAML interface. Protect your infra with an autonomous sentinel.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <button className="w-full sm:w-auto bg-blue-600 hover:bg-blue-500 text-white px-8 py-4 rounded-xl font-bold shadow-xl shadow-blue-500/25 flex items-center justify-center gap-2 group transition-all">
                Get Started
                <ChevronRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </button>
              <button className="w-full sm:w-auto bg-zinc-900 border border-zinc-800 hover:bg-zinc-800 text-white px-8 py-4 rounded-xl font-bold flex items-center justify-center gap-2 transition-all">
                <Terminal className="w-5 h-5 text-zinc-500" />
                View Documentation
              </button>
            </div>
          </motion.div>
        </div>

        <motion.div 
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.3, duration: 0.8 }}
          className="mt-20 max-w-5xl mx-auto glass rounded-2xl overflow-hidden shadow-2xl border-zinc-800/50"
        >
          <div className="bg-zinc-900/50 px-4 py-3 border-b border-zinc-800 flex items-center justify-between">
            <div className="flex gap-2">
              <div className="w-3 h-3 rounded-full bg-red-500/20 border border-red-500/40" />
              <div className="w-3 h-3 rounded-full bg-yellow-500/20 border border-yellow-500/40" />
              <div className="w-3 h-3 rounded-full bg-green-500/20 border border-green-500/40" />
            </div>
            <span className="text-[10px] text-zinc-500 font-mono tracking-widest uppercase">policy.yaml</span>
            <div className="w-10" />
          </div>
          <div className="p-6 md:p-8 font-mono text-sm md:text-base bg-zinc-950/50">
            <pre className="text-zinc-300 leading-relaxed overflow-x-auto">
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

const FeatureCard = ({ icon: Icon, title, description, delay }: { icon: any, title: string, description: string, delay: number }) => (
  <motion.div 
    initial={{ opacity: 0, y: 20 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true }}
    transition={{ delay, duration: 0.5 }}
    className="p-8 rounded-2xl bg-zinc-900/30 border border-zinc-800 hover:border-blue-500/30 hover:bg-blue-500/[0.02] transition-all group"
  >
    <div className="w-12 h-12 bg-zinc-800 rounded-xl flex items-center justify-center mb-6 group-hover:bg-blue-600 group-hover:text-white transition-colors">
      <Icon className="w-6 h-6 text-blue-500 group-hover:text-white transition-colors" />
    </div>
    <h3 className="text-xl font-bold mb-3 text-white">{title}</h3>
    <p className="text-zinc-400 leading-relaxed">{description}</p>
  </motion.div>
)

const Features = () => {
  return (
    <section id="features" className="py-24 bg-zinc-950">
      <div className="container mx-auto px-6">
        <div className="text-center max-w-2xl mx-auto mb-20">
          <h2 className="text-3xl md:text-5xl font-bold mb-6 tracking-tight">Powerful by Default.</h2>
          <p className="text-zinc-400 text-lg leading-relaxed">
            Rampart combines the raw performance of Linux kernel technologies with the ease of use of modern cloud APIs.
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <FeatureCard 
            icon={Cpu}
            title="eBPF/XDP Engine"
            description="Ultra-high performance packet filtering at the driver level. Millions of packets per second with microsecond latency."
            delay={0.1}
          />
          <FeatureCard 
            icon={Bot}
            title="Autonomous Sentinel"
            description="Built-in IPS module that analyzes traffic patterns in real-time and automatically locks down nodes during attacks."
            delay={0.2}
          />
          <FeatureCard 
            icon={Cloud}
            title="Unified Multi-Cloud"
            description="Orchestrate security groups across AWS, GCP, and Azure simultaneously using a single YAML policy set."
            delay={0.3}
          />
          <FeatureCard 
            icon={RefreshCw}
            title="Self-Healing"
            description="Background watchdog service continuously monitors and corrects firewall state to ensure 100% compliance."
            delay={0.4}
          />
          <FeatureCard 
            icon={Search}
            title="DPI Awareness"
            description="Deep packet inspection for Layer-7 filtering. Block malicious DNS queries and HTTP hosts natively."
            delay={0.5}
          />
          <FeatureCard 
            icon={Activity}
            title="Enterprise Monitoring"
            description="Built-in Prometheus metrics, live SSE events, and tamper-evident audit logs with cryptographic hash-chains."
            delay={0.6}
          />
        </div>
      </div>
    </section>
  )
}

const Architecture = () => {
  return (
    <section id="architecture" className="py-24 bg-zinc-950 relative overflow-hidden">
      <div className="container mx-auto px-6">
        <div className="flex flex-col lg:flex-row items-center gap-16">
          <div className="lg:w-1/2">
            <h2 className="text-4xl font-bold mb-8 tracking-tight">The Distributed <br /><span className="text-blue-500">Control Plane.</span></h2>
            <div className="space-y-6">
              {[
                { title: 'Raft-based Consensus', desc: 'Secure, distributed policy synchronization with mTLS encryption across all nodes.' },
                { title: 'Backend Abstraction Layer', desc: 'Plug-and-play support for nftables, iptables, eBPF, and Cloud Security Groups.' },
                { title: 'AI-Ready Interface', desc: 'Native Model Context Protocol (MCP) support for seamless orchestration by AI agents.' }
              ].map((item, i) => (
                <div key={i} className="flex gap-4">
                  <div className="mt-1">
                    <CheckCircle2 className="w-6 h-6 text-blue-500" />
                  </div>
                  <div>
                    <h4 className="text-lg font-bold text-white mb-1">{item.title}</h4>
                    <p className="text-zinc-400">{item.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
          <div className="lg:w-1/2 relative">
            <div className="relative z-10 glass p-8 rounded-3xl border-zinc-800">
               <div className="grid grid-cols-2 gap-4">
                  <div className="p-4 rounded-xl bg-blue-600/10 border border-blue-500/20 text-center">
                    <Cloud className="w-8 h-8 text-blue-500 mx-auto mb-2" />
                    <span className="text-xs font-bold text-blue-400 uppercase tracking-widest">Cloud API</span>
                  </div>
                  <div className="p-4 rounded-xl bg-indigo-600/10 border border-indigo-500/20 text-center">
                    <Cpu className="w-8 h-8 text-indigo-500 mx-auto mb-2" />
                    <span className="text-xs font-bold text-indigo-400 uppercase tracking-widest">eBPF Engine</span>
                  </div>
                  <div className="col-span-2 p-6 rounded-xl bg-zinc-800/30 border border-zinc-700 text-center">
                    <Lock className="w-8 h-8 text-zinc-400 mx-auto mb-2" />
                    <span className="text-sm font-bold text-white uppercase tracking-widest italic">Raft Consensus Core</span>
                  </div>
                  <div className="p-4 rounded-xl bg-emerald-600/10 border border-emerald-500/20 text-center">
                    <Shield className="w-8 h-8 text-emerald-500 mx-auto mb-2" />
                    <span className="text-xs font-bold text-emerald-400 uppercase tracking-widest">nftables</span>
                  </div>
                  <div className="p-4 rounded-xl bg-orange-600/10 border border-orange-500/20 text-center">
                    <Activity className="w-8 h-8 text-orange-500 mx-auto mb-2" />
                    <span className="text-xs font-bold text-orange-400 uppercase tracking-widest">SIEM / Log</span>
                  </div>
               </div>
            </div>
            {/* Background Blobs */}
            <div className="absolute -top-10 -right-10 w-64 h-64 bg-blue-600/20 blur-3xl rounded-full" />
            <div className="absolute -bottom-10 -left-10 w-64 h-64 bg-indigo-600/10 blur-3xl rounded-full" />
          </div>
        </div>
      </div>
    </section>
  )
}

const Footer = () => (
  <footer className="py-12 border-t border-zinc-900 bg-zinc-950">
    <div className="container mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-8">
      <div className="flex items-center gap-2">
        <Shield className="text-blue-500 w-6 h-6" />
        <span className="text-lg font-bold tracking-tight text-white">RAMPART</span>
      </div>
      <p className="text-zinc-500 text-sm">
        Licensed under Apache License 2.0. Built for the modern cloud.
      </p>
      <div className="flex items-center gap-6">
        <a href="https://github.com/ersinkoc/Rampart" className="text-zinc-500 hover:text-white transition-colors">
          <Github className="w-5 h-5" />
        </a>
        <a href="#" className="text-zinc-500 hover:text-white transition-colors">
          <Globe className="w-5 h-5" />
        </a>
      </div>
    </div>
  </footer>
)

const App = () => {
  return (
    <main className="bg-zinc-950 min-h-screen">
      <Navbar />
      <Hero />
      <Features />
      <Architecture />
      <section className="py-24 container mx-auto px-6 text-center">
        <div className="glass p-12 md:p-20 rounded-[3rem] border-zinc-800 relative overflow-hidden">
          <h2 className="text-4xl md:text-6xl font-extrabold mb-8 relative z-10 tracking-tight">Ready to secure your <br /><span className="text-blue-500">entire fleet?</span></h2>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4 relative z-10">
            <a 
              href="https://github.com/ersinkoc/Rampart"
              className="w-full sm:w-auto bg-white text-zinc-950 px-10 py-5 rounded-2xl font-bold hover:bg-zinc-200 transition-all text-lg shadow-2xl"
            >
              Get Started Now
            </a>
          </div>
          <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-blue-600/10 via-transparent to-indigo-600/5" />
        </div>
      </section>
      <Footer />
    </main>
  )
}

export default App
