import React, { useState, useEffect } from 'react';
import { 
  Shield, 
  Bot, 
  Activity, 
  Server, 
  Zap, 
  AlertTriangle, 
  Lock, 
  Globe, 
  Cpu,
  RefreshCw,
  Search,
  CheckCircle2
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const StatCard = ({ icon: Icon, title, value, subValue, trend, color }: any) => (
  <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 p-6 rounded-2xl shadow-sm hover:shadow-md transition-all">
    <div className="flex items-center justify-between mb-4">
      <div className={`p-2 rounded-xl bg-${color}-500/10 text-${color}-600 dark:text-${color}-400`}>
        <Icon className="w-5 h-5" />
      </div>
      {trend && (
        <span className={`text-xs font-bold ${trend > 0 ? 'text-emerald-500' : 'text-rose-500'}`}>
          {trend > 0 ? '+' : ''}{trend}%
        </span>
      )}
    </div>
    <h3 className="text-zinc-500 dark:text-zinc-400 text-sm font-medium">{title}</h3>
    <div className="mt-1 flex items-baseline gap-2">
      <p className="text-2xl font-bold dark:text-white">{value}</p>
      {subValue && <span className="text-xs text-zinc-400 font-mono">{subValue}</span>}
    </div>
  </div>
);

const Dashboard: React.FC = () => {
  const [events, setEvents] = useState<any[]>([]);
  const [threatScore, setThreatScore] = useState(12);

  // Simulate live threat updates
  useEffect(() => {
    const interval = setInterval(() => {
      setThreatScore(prev => Math.max(0, Math.min(100, prev + (Math.random() > 0.7 ? 5 : -2))));
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight dark:text-white">Autonomous Sentinel</h1>
          <p className="text-zinc-500 mt-1">Real-time network defense and cluster health monitoring.</p>
        </div>
        <div className="flex items-center gap-3 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 px-4 py-2 rounded-full border border-emerald-500/20">
          <Activity className="w-4 h-4 animate-pulse" />
          <span className="text-sm font-bold uppercase tracking-wider">Live & Protected</span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
        <StatCard 
          icon={Shield} 
          title="Active Rules" 
          value="156" 
          subValue="across 4 backends" 
          trend={12} 
          color="blue" 
        />
        <StatCard 
          icon={Bot} 
          title="Threat Score" 
          value={threatScore} 
          subValue="Sentinel active" 
          trend={-5} 
          color="indigo" 
        />
        <StatCard 
          icon={Server} 
          title="Cluster Nodes" 
          value="5" 
          subValue="Quorum achieved" 
          color="emerald" 
        />
        <StatCard 
          icon={Zap} 
          title="eBPF Packets" 
          value="1.2M" 
          subValue="last 5 min" 
          trend={8} 
          color="orange" 
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Feed */}
        <div className="lg:col-span-2 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl overflow-hidden shadow-sm">
          <div className="px-8 py-6 border-b border-zinc-100 dark:border-zinc-800 flex items-center justify-between">
            <h3 className="text-lg font-bold dark:text-white flex items-center gap-2">
              <Activity className="w-5 h-5 text-blue-500" />
              Live Security Events
            </h3>
            <button className="text-sm text-zinc-400 hover:text-white flex items-center gap-2">
              <RefreshCw className="w-4 h-4" />
              Refresh
            </button>
          </div>
          <div className="divide-y divide-zinc-50 dark:divide-zinc-800/50">
            {[
              { id: '1', type: 'ips-ban', title: 'IP Quarantine Triggered', actor: 'Autonomous Sentinel', resource: '192.168.1.45', time: '2 min ago', severity: 'high' },
              { id: '2', type: 'policy-apply', title: 'Global Ingress Updated', actor: 'ersin', resource: 'rampartfw.com/v1', time: '14 min ago', severity: 'info' },
              { id: '3', type: 'dpi-signal', title: 'DPI Anomaly Detected', actor: 'DPI Engine', resource: 'malicious-domain.com', time: '1 hour ago', severity: 'medium' },
              { id: '4', type: 'cluster-sync', title: 'Raft State Synchronized', actor: 'node-3', resource: 'Term 42', time: '2 hours ago', severity: 'info' },
            ].map(event => (
              <div key={event.id} className="px-8 py-5 hover:bg-zinc-50 dark:hover:bg-white/[0.02] transition-colors flex items-center justify-between group">
                <div className="flex items-center gap-5">
                  <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${
                    event.severity === 'high' ? 'bg-rose-500/10 text-rose-500' : 
                    event.severity === 'medium' ? 'bg-amber-500/10 text-amber-500' : 
                    'bg-blue-500/10 text-blue-500'
                  }`}>
                    {event.type === 'ips-ban' ? <Lock className="w-5 h-5" /> : 
                     event.type === 'dpi-signal' ? <Search className="w-5 h-5" /> : 
                     <Zap className="w-5 h-5" />}
                  </div>
                  <div>
                    <p className="text-sm font-bold dark:text-white">{event.title}</p>
                    <p className="text-xs text-zinc-500 mt-0.5">{event.actor} • {event.resource}</p>
                  </div>
                </div>
                <span className="text-xs text-zinc-400 font-medium group-hover:text-zinc-300 transition-colors">{event.time}</span>
              </div>
            ))}
          </div>
          <div className="p-4 bg-zinc-50 dark:bg-zinc-950/50 text-center border-t border-zinc-100 dark:border-zinc-800">
            <button className="text-sm font-bold text-blue-500 hover:text-blue-400">View Full Audit Trail</button>
          </div>
        </div>

        {/* Sidebar Status */}
        <div className="space-y-6">
          <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 p-8 rounded-3xl shadow-sm relative overflow-hidden">
             <h3 className="text-lg font-bold dark:text-white mb-6 flex items-center gap-2">
               <Cpu className="w-5 h-5 text-indigo-500" />
               Engine Orchestration
             </h3>
             <div className="space-y-5 relative z-10">
                {[
                  { name: 'nftables', status: 'Active', color: 'emerald' },
                  { name: 'eBPF/XDP', status: 'Fast Path', color: 'blue' },
                  { name: 'AWS Cloud', status: 'Synchronized', color: 'orange' },
                  { name: 'Raft consensus', status: 'Leader', color: 'purple' },
                ].map(item => (
                  <div key={item.name} className="flex justify-between items-center">
                    <span className="text-zinc-400 text-sm font-medium">{item.name}</span>
                    <span className={`text-[10px] font-extrabold uppercase tracking-widest px-2 py-0.5 rounded bg-${item.color}-500/10 text-${item.color}-500 border border-${item.color}-500/20`}>
                      {item.status}
                    </span>
                  </div>
                ))}
             </div>
             <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-500/5 blur-3xl rounded-full" />
          </div>

          <div className="bg-gradient-to-br from-blue-600 to-indigo-700 p-8 rounded-3xl shadow-lg text-white">
            <h3 className="text-lg font-bold mb-2">Need AI assistance?</h3>
            <p className="text-blue-100 text-sm mb-6 leading-relaxed">Use the built-in MCP server to let AI agents manage your policy fleet autonomously.</p>
            <button className="w-full bg-white text-blue-600 py-3 rounded-xl font-bold hover:bg-zinc-100 transition-colors flex items-center justify-center gap-2 shadow-xl">
              <Bot className="w-5 h-5" />
              Open AI Console
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
