import React from 'react';
import { 
  Network, 
  Server, 
  Activity, 
  Shield, 
  Cpu, 
  RefreshCw, 
  CheckCircle2, 
  Lock,
  Globe,
  ArrowUpRight,
  Database
} from 'lucide-react';
import { motion } from 'framer-motion';

const Cluster: React.FC = () => {
  const nodes = [
    { id: 'node-1', state: 'leader', backend: 'eBPF+nftables', rules: 156, latency: '0.02ms', health: 'healthy', location: 'Frankfurt' },
    { id: 'node-2', state: 'follower', backend: 'nftables', rules: 156, latency: '12ms', health: 'healthy', location: 'London' },
    { id: 'node-3', state: 'follower', backend: 'eBPF', rules: 156, latency: '45ms', health: 'healthy', location: 'New York' },
    { id: 'node-4', state: 'follower', backend: 'iptables', rules: 156, latency: '8ms', health: 'degraded', location: 'Singapore' },
    { id: 'node-5', state: 'candidate', backend: 'AWS SG', rules: 156, latency: '120ms', health: 'healthy', location: 'Tokyo' },
  ];

  return (
    <div className="space-y-10 animate-in fade-in duration-700">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div>
          <h1 className="text-4xl font-black tracking-tight dark:text-white text-zinc-950 flex items-center gap-4">
            <Network className="w-10 h-10 text-blue-600" />
            Distributed Mesh
          </h1>
          <p className="text-zinc-500 font-medium text-lg mt-1">Global state orchestration via Raft consensus (mTLS).</p>
        </div>
        <div className="flex items-center gap-3 bg-blue-600/10 text-blue-600 px-5 py-2.5 rounded-2xl border border-blue-500/20">
          <Lock className="w-4 h-4" />
          <span className="text-xs font-black uppercase tracking-[0.2em]">Quorum Achieved (4/5)</span>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Consensus Stats */}
        <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-6 h-fit">
           <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 p-8 rounded-[3rem] shadow-sm relative overflow-hidden group">
              <h3 className="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-6 flex items-center gap-2">
                 <Database className="w-4 h-4" />
                 Raft Integrity
              </h3>
              <div className="space-y-4">
                 <div className="flex justify-between items-end">
                    <span className="text-3xl font-black dark:text-white">Term 42</span>
                    <span className="text-[10px] font-bold text-emerald-500 uppercase">Synchronized</span>
                 </div>
                 <div className="w-full h-1.5 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden">
                    <div className="h-full bg-blue-600 w-full" />
                 </div>
                 <div className="flex justify-between text-[10px] font-black text-zinc-400 uppercase tracking-tighter">
                    <span>Index: 12,841</span>
                    <span>Commit: 12,841</span>
                 </div>
              </div>
           </div>

           <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 p-8 rounded-[3rem] shadow-sm relative overflow-hidden group">
              <h3 className="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-6 flex items-center gap-2">
                 <Activity className="w-4 h-4" />
                 Cluster Load
              </h3>
              <div className="space-y-4">
                 <div className="flex justify-between items-end">
                    <span className="text-3xl font-black dark:text-white">8.4k</span>
                    <span className="text-[10px] font-bold text-zinc-400 uppercase">RPC / SEC</span>
                 </div>
                 <div className="w-full h-1.5 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden">
                    <div className="h-full bg-indigo-500 w-1/3" />
                 </div>
                 <div className="flex justify-between text-[10px] font-black text-zinc-400 uppercase tracking-tighter">
                    <span>Avg Jitter: 0.12ms</span>
                    <span>MTU: 1500</span>
                 </div>
              </div>
           </div>
        </div>

        {/* Global Topology Visualizer Placeholder */}
        <div className="lg:col-span-1 bg-zinc-950 rounded-[3rem] border border-white/5 p-1 relative overflow-hidden h-full min-h-[250px]">
           <div className="absolute inset-0 opacity-20 pointer-events-none">
              <div className="w-full h-full bg-[radial-gradient(circle_at_center,rgba(37,99,235,0.4),transparent_70%)]" />
           </div>
           <div className="relative h-full w-full flex flex-col items-center justify-center p-8 text-center">
              <div className="w-16 h-16 rounded-full border-2 border-blue-500/50 border-t-blue-500 animate-spin-slow mb-4 flex items-center justify-center">
                 <Globe className="w-8 h-8 text-blue-500" />
              </div>
              <h4 className="text-sm font-black text-white uppercase tracking-[0.3em]">Mesh Topology</h4>
              <p className="text-[10px] text-zinc-500 font-bold mt-2 italic">Real-time graph visualization...</p>
           </div>
        </div>
      </div>

      <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] overflow-hidden shadow-sm">
        <div className="px-10 py-8 border-b border-zinc-100 dark:border-white/5 flex items-center justify-between">
           <h3 className="text-xl font-black dark:text-white text-zinc-950">Active Peer Manifest</h3>
           <div className="flex gap-4">
              <div className="flex items-center gap-2">
                 <div className="w-2 h-2 rounded-full bg-emerald-500" />
                 <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">Leader</span>
              </div>
              <div className="flex items-center gap-2">
                 <div className="w-2 h-2 rounded-full bg-blue-500" />
                 <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">Follower</span>
              </div>
           </div>
        </div>
        
        <div className="overflow-x-auto">
          <table className="w-full text-left font-medium">
            <thead>
              <tr className="bg-zinc-50/50 dark:bg-white/[0.02] text-[10px] font-black text-zinc-500 uppercase tracking-[0.3em] border-b border-zinc-100 dark:border-white/5">
                <th className="px-10 py-5">Node Identity</th>
                <th className="px-10 py-5">Role</th>
                <th className="px-10 py-5">Orchestrator Backend</th>
                <th className="px-10 py-5">Managed Rules</th>
                <th className="px-10 py-5 text-right">Heartbeat</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-50 dark:divide-white/5">
              {nodes.map((node) => (
                <tr key={node.id} className="group hover:bg-zinc-50 dark:hover:bg-white/[0.01] transition-colors">
                  <td className="px-10 py-6">
                    <div className="flex items-center gap-4">
                       <div className={cn(
                         "w-12 h-12 rounded-2xl flex items-center justify-center transition-all duration-500 group-hover:scale-110",
                         node.state === 'leader' ? "bg-emerald-500/10 text-emerald-500 shadow-xl shadow-emerald-500/10" : "bg-zinc-100 dark:bg-zinc-800 text-zinc-400"
                       )}>
                         <Server className="w-6 h-6" />
                       </div>
                       <div>
                          <p className="font-bold dark:text-white text-zinc-950 underline decoration-zinc-300 dark:decoration-zinc-800 underline-offset-4">{node.id}</p>
                          <p className="text-[10px] text-zinc-500 font-black mt-1 flex items-center gap-1">
                             <Globe className="w-3 h-3" />
                             {node.location}
                          </p>
                       </div>
                    </div>
                  </td>
                  <td className="px-10 py-6">
                     <span className={cn(
                       "px-4 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest border",
                       node.state === 'leader' ? "bg-emerald-500/10 text-emerald-600 border-emerald-500/20" : 
                       node.state === 'candidate' ? "bg-amber-500/10 text-amber-600 border-amber-500/20" :
                       "bg-blue-500/10 text-blue-600 border-blue-500/20"
                     )}>
                        {node.state}
                     </span>
                  </td>
                  <td className="px-10 py-6 text-sm font-bold dark:text-zinc-300 text-zinc-600 italic">
                     {node.backend}
                  </td>
                  <td className="px-10 py-6">
                     <div className="flex items-center gap-3">
                        <span className="text-lg font-black dark:text-white">{node.rules}</span>
                        <div className="h-4 w-[1px] bg-zinc-200 dark:bg-zinc-800" />
                        <CheckCircle2 className="w-4 h-4 text-emerald-500" />
                     </div>
                  </td>
                  <td className="px-10 py-6 text-right">
                     <div className="inline-flex items-center gap-2 bg-zinc-50 dark:bg-zinc-800/50 px-3 py-1 rounded-lg">
                        <Activity className="w-3 h-3 text-zinc-400" />
                        <span className="text-[10px] font-black font-mono text-zinc-500">{node.latency}</span>
                     </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Cluster;

function cn(...inputs: any[]) {
  return inputs.filter(Boolean).join(' ');
}
