import React, { useState } from 'react';
import { 
  Zap, 
  Send, 
  Search, 
  ArrowRight, 
  CheckCircle2, 
  XCircle, 
  AlertCircle,
  Clock,
  Globe,
  Terminal,
  Cpu
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const Simulator: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);

  const runSimulation = () => {
    setLoading(true);
    // Mock simulation delay
    setTimeout(() => {
      setResult({
        verdict: 'ACCEPT',
        rule: 'allow-http',
        path: 'src 10.0.1.50 ∈ 0.0.0.0/0, proto tcp, dport 80',
        evaluated: 12,
        duration: '0.45ms'
      });
      setLoading(false);
    }, 800);
  };

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight dark:text-white">Traffic Simulator</h1>
          <p className="text-zinc-500 mt-1">Predictive analysis: Test packets against the current policy engine.</p>
        </div>
        <div className="flex items-center gap-3 bg-blue-500/10 text-blue-600 dark:text-blue-400 px-4 py-2 rounded-full border border-blue-500/20">
          <Cpu className="w-4 h-4" />
          <span className="text-sm font-bold uppercase tracking-wider text-xs">Simulating v0.1.0</span>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-[2.5rem] p-8 shadow-sm">
          <h3 className="text-lg font-bold dark:text-white mb-8 flex items-center gap-2">
            <Terminal className="w-5 h-5 text-zinc-400" />
            Packet Properties
          </h3>
          
          <div className="space-y-6">
            <div className="grid grid-cols-2 gap-6">
              <div className="space-y-2">
                <label className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest ml-1">Source IP</label>
                <input type="text" defaultValue="10.0.1.50" className="w-full px-4 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 focus:ring-2 focus:ring-blue-500/50 outline-none text-sm font-mono" />
              </div>
              <div className="space-y-2">
                <label className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest ml-1">Destination IP</label>
                <input type="text" defaultValue="192.168.1.10" className="w-full px-4 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 focus:ring-2 focus:ring-blue-500/50 outline-none text-sm font-mono" />
              </div>
            </div>

            <div className="grid grid-cols-3 gap-6">
              <div className="space-y-2">
                <label className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest ml-1">Protocol</label>
                <select className="w-full px-4 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 focus:ring-2 focus:ring-blue-500/50 outline-none text-sm font-bold">
                  <option>TCP</option>
                  <option>UDP</option>
                  <option>ICMP</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest ml-1">Dest Port</label>
                <input type="number" defaultValue="80" className="w-full px-4 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 focus:ring-2 focus:ring-blue-500/50 outline-none text-sm font-mono" />
              </div>
              <div className="space-y-2">
                <label className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest ml-1">Direction</label>
                <select className="w-full px-4 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 focus:ring-2 focus:ring-blue-500/50 outline-none text-sm font-bold">
                  <option>Inbound</option>
                  <option>Outbound</option>
                </select>
              </div>
            </div>

            <div className="pt-4">
              <button 
                onClick={runSimulation}
                disabled={loading}
                className="w-full py-4 rounded-2xl bg-blue-600 hover:bg-blue-500 text-white font-bold shadow-xl shadow-blue-500/25 flex items-center justify-center gap-2 transition-all disabled:opacity-50"
              >
                {loading ? <RefreshCw className="w-5 h-5 animate-spin" /> : <Send className="w-5 h-5" />}
                {loading ? 'Analyzing Path...' : 'Inject Simulated Packet'}
              </button>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-6">
          <AnimatePresence mode="wait">
            {!result && !loading ? (
              <motion.div 
                key="empty"
                initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
                className="flex-1 border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-[2.5rem] flex flex-col items-center justify-center p-12 text-center"
              >
                <div className="w-20 h-20 rounded-full bg-zinc-100 dark:bg-white/[0.02] flex items-center justify-center mb-6">
                  <Search className="w-10 h-10 text-zinc-300 dark:text-zinc-700" />
                </div>
                <h4 className="text-xl font-bold dark:text-white mb-2">No Simulation Data</h4>
                <p className="text-sm text-zinc-500 max-w-xs">Configure the packet properties and click simulate to see the evaluation path.</p>
              </motion.div>
            ) : loading ? (
               <motion.div 
                 key="loading"
                 initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
                 className="flex-1 glass rounded-[2.5rem] flex flex-col items-center justify-center p-12 overflow-hidden relative"
               >
                  <div className="w-24 h-24 rounded-full border-4 border-blue-500/20 border-t-blue-500 animate-spin mb-8" />
                  <p className="text-lg font-bold animate-pulse text-blue-500">Traversing Kernel Tables...</p>
                  <div className="absolute inset-0 bg-blue-500/5 pointer-events-none" />
               </motion.div>
            ) : (
              <motion.div 
                key="result"
                initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}
                className="flex-1 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-[2.5rem] overflow-hidden flex flex-col shadow-sm"
              >
                <div className={cn(
                  "p-8 text-center",
                  result.verdict === 'ACCEPT' ? 'bg-emerald-500/10' : 'bg-rose-500/10'
                )}>
                  <div className={cn(
                    "w-16 h-16 rounded-2xl mx-auto mb-4 flex items-center justify-center shadow-lg",
                    result.verdict === 'ACCEPT' ? 'bg-emerald-500 text-white shadow-emerald-500/20' : 'bg-rose-500 text-white shadow-rose-500/20'
                  )}>
                    {result.verdict === 'ACCEPT' ? <CheckCircle2 className="w-8 h-8" /> : <XCircle className="w-8 h-8" />}
                  </div>
                  <h2 className={cn(
                    "text-3xl font-black tracking-tighter",
                    result.verdict === 'ACCEPT' ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400'
                  )}>
                    {result.verdict}
                  </h2>
                  <p className="text-zinc-500 text-sm mt-1 font-medium">Policy evaluation complete.</p>
                </div>
                
                <div className="p-8 flex-1 space-y-8">
                   <div className="grid grid-cols-2 gap-4">
                      <div className="p-4 rounded-2xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-100 dark:border-zinc-800">
                         <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest mb-1">Evaluated</p>
                         <p className="text-sm font-bold dark:text-white">{result.evaluated} rules</p>
                      </div>
                      <div className="p-4 rounded-2xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-100 dark:border-zinc-800">
                         <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest mb-1">Duration</p>
                         <p className="text-sm font-bold dark:text-white">{result.duration}</p>
                      </div>
                   </div>

                   <div className="space-y-4">
                      <h4 className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest flex items-center gap-2">
                        <ArrowRight className="w-3 h-3" />
                        Match Trace
                      </h4>
                      <div className="p-5 rounded-2xl bg-zinc-900 border border-zinc-800 font-mono text-xs text-blue-400 leading-relaxed shadow-inner">
                         {result.path}
                      </div>
                      <div className="flex items-center gap-3 p-4 rounded-xl bg-blue-500/5 border border-blue-500/10">
                         <Info className="w-4 h-4 text-blue-500 shrink-0" />
                         <p className="text-xs text-blue-600 dark:text-blue-300">Packet matched rule <span className="font-bold underline">{result.rule}</span> at priority <span className="font-bold">500</span>.</p>
                      </div>
                   </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
};

export default Simulator;

function cn(...inputs: any[]) {
  return inputs.filter(Boolean).join(' ');
}
