import React, { useState, useEffect } from 'react';
import { 
  Shield, 
  Bot, 
  Activity, 
  Server, 
  Zap, 
  Lock, 
  Cpu,
  RefreshCw,
  Search,
  CheckCircle2,
  TrendingUp,
  AlertOctagon,
  MousePointer2
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  AreaChart, 
  Area, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from 'recharts';

const data = [
  { time: '10:00', packets: 400, threats: 24 },
  { time: '10:05', packets: 300, threats: 18 },
  { time: '10:10', packets: 600, threats: 90 }, // Spike
  { time: '10:15', packets: 800, threats: 40 },
  { time: '10:20', packets: 500, threats: 30 },
  { time: '10:25', packets: 900, threats: 15 },
  { time: '10:30', packets: 1100, threats: 10 },
];

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444'];

const StatCard = ({ icon: Icon, title, value, subValue, trend, color }: any) => (
  <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 p-8 rounded-[2.5rem] shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-500 group relative overflow-hidden">
    <div className="flex items-center justify-between mb-6 relative z-10">
      <div className={cn(
        "w-14 h-14 rounded-2xl flex items-center justify-center transition-all duration-500",
        `bg-${color}-500/10 text-${color}-600 dark:text-${color}-400 group-hover:bg-${color}-500 group-hover:text-white`
      )}>
        <Icon className="w-7 h-7" />
      </div>
      {trend && (
        <div className={cn(
          "px-3 py-1 rounded-full text-xs font-black tracking-tight",
          trend > 0 ? "bg-emerald-500/10 text-emerald-600" : "bg-rose-500/10 text-rose-600"
        )}>
          {trend > 0 ? '↑' : '↓'} {Math.abs(trend)}%
        </div>
      )}
    </div>
    <div className="relative z-10">
        <h3 className="text-zinc-500 dark:text-zinc-400 text-sm font-bold uppercase tracking-widest">{title}</h3>
        <div className="mt-2 flex items-baseline gap-3">
          <p className="text-4xl font-black dark:text-white text-zinc-950 tracking-tighter">{value}</p>
          {subValue && <span className="text-xs text-zinc-400 font-bold uppercase tracking-tighter">{subValue}</span>}
        </div>
    </div>
    <div className={cn("absolute -bottom-10 -right-10 w-32 h-32 blur-3xl rounded-full opacity-0 group-hover:opacity-20 transition-opacity duration-700", `bg-${color}-500`)} />
  </div>
);

const Dashboard: React.FC = () => {
  const [threatScore, setThreatScore] = useState(12);

  useEffect(() => {
    const interval = setInterval(() => {
      setThreatScore(prev => Math.max(0, Math.min(100, prev + (Math.random() > 0.8 ? 8 : -3))));
    }, 4000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-10 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
        <div>
          <h1 className="text-4xl md:text-5xl font-black tracking-tight dark:text-white text-zinc-950">
            Sentinel <span className="text-blue-600 italic">Cockpit</span>
          </h1>
          <p className="text-zinc-500 font-medium text-lg mt-1">Autonomous monitoring for node-1 (Cluster Leader)</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex -space-x-3">
             {[1, 2, 3, 4, 5].map(n => (
               <div key={n} className="w-10 h-10 rounded-full border-4 border-white dark:border-zinc-950 bg-zinc-800 flex items-center justify-center text-[10px] font-black text-white">
                  N{n}
               </div>
             ))}
          </div>
          <div className="h-10 w-[1px] bg-zinc-200 dark:bg-zinc-800 mx-2" />
          <div className="flex items-center gap-3 bg-emerald-500/10 text-emerald-600 px-5 py-2.5 rounded-2xl border border-emerald-500/20">
            <div className="w-2 h-2 rounded-full bg-emerald-500 animate-ping" />
            <span className="text-xs font-black uppercase tracking-[0.2em]">Operational</span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard icon={Shield} title="Managed Rules" value="1,248" subValue="Across all planes" trend={14} color="blue" />
        <StatCard icon={Bot} title="Sentinel Risk" value={`${threatScore}%`} subValue="Dynamic Score" trend={-4} color="indigo" />
        <StatCard icon={Activity} title="Avg Latency" value="0.04" subValue="Microseconds" trend={-2} color="emerald" />
        <StatCard icon={Zap} title="Blocked Attacks" value="842" subValue="Last 24 hours" trend={28} color="rose" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] p-10 shadow-sm relative overflow-hidden group">
          <div className="flex items-center justify-between mb-10">
             <div>
                <h3 className="text-2xl font-black dark:text-white text-zinc-950 tracking-tighter">Traffic Velocity</h3>
                <p className="text-sm text-zinc-500 font-medium">Real-time throughput and mitigation trends</p>
             </div>
             <div className="flex gap-4">
                <div className="flex items-center gap-2">
                   <div className="w-3 h-3 rounded-full bg-blue-500" />
                   <span className="text-xs font-bold text-zinc-400">Packets</span>
                </div>
                <div className="flex items-center gap-2">
                   <div className="w-3 h-3 rounded-full bg-rose-500" />
                   <span className="text-xs font-bold text-zinc-400">Threats</span>
                </div>
             </div>
          </div>
          
          <div className="h-[350px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={data}>
                <defs>
                  <linearGradient id="colorPkts" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorThreats" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#ef4444" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#ef4444" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <XAxis dataKey="time" axisLine={false} tickLine={false} tick={{fill: '#71717a', fontSize: 10, fontWeight: 'bold'}} />
                <YAxis hide />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#09090b', border: '1px solid #27272a', borderRadius: '16px', fontSize: '12px', fontWeight: 'bold' }}
                  itemStyle={{ color: '#fff' }}
                />
                <Area type="monotone" dataKey="packets" stroke="#3b82f6" strokeWidth={4} fillOpacity={1} fill="url(#colorPkts)" />
                <Area type="monotone" dataKey="threats" stroke="#ef4444" strokeWidth={4} fillOpacity={1} fill="url(#colorThreats)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
          <div className="absolute top-0 right-0 w-64 h-64 bg-blue-500/5 blur-[100px] pointer-events-none" />
        </div>

        <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] p-10 flex flex-col shadow-sm">
           <h3 className="text-2xl font-black dark:text-white text-zinc-950 tracking-tighter mb-2 text-center">Threat Vector</h3>
           <p className="text-sm text-zinc-500 font-bold uppercase tracking-widest text-center mb-8">Classification</p>
           
           <div className="flex-1 min-h-[250px]">
             <ResponsiveContainer width="100%" height="100%">
               <PieChart>
                 <Pie
                   data={[
                     { name: 'L3/L4 Flood', value: 45 },
                     { name: 'DPI Anomaly', value: 30 },
                     { name: 'K8s Violation', value: 15 },
                     { name: 'Brute Force', value: 10 },
                   ]}
                   innerRadius={60}
                   outerRadius={100}
                   paddingAngle={8}
                   dataKey="value"
                 >
                   {[0,1,2,3].map((entry, index) => (
                     <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                   ))}
                 </Pie>
                 <Tooltip />
               </PieChart>
             </ResponsiveContainer>
           </div>
           
           <div className="space-y-4 mt-6">
              {[
                { name: 'L3/L4 Flood', color: 'bg-blue-500' },
                { name: 'DPI Anomaly', color: 'bg-emerald-500' },
                { name: 'K8s Violation', color: 'bg-amber-500' },
                { name: 'Brute Force', color: 'bg-rose-500' },
              ].map(item => (
                <div key={item.name} className="flex items-center justify-between">
                   <div className="flex items-center gap-3">
                      <div className={cn("w-2 h-2 rounded-full", item.color)} />
                      <span className="text-xs font-bold text-zinc-500">{item.name}</span>
                   </div>
                   <span className="text-xs font-mono dark:text-zinc-300">ACTIVE</span>
                </div>
              ))}
           </div>
        </div>
      </div>

      {/* Modern Event Table */}
      <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] overflow-hidden shadow-sm">
        <div className="px-10 py-8 border-b border-zinc-100 dark:border-white/5 flex items-center justify-between bg-zinc-50/50 dark:bg-white/[0.01]">
          <h3 className="text-xl font-black dark:text-white text-zinc-950 flex items-center gap-3">
            <AlertOctagon className="w-6 h-6 text-rose-500" />
            Immediate Threats
          </h3>
          <button className="px-5 py-2 rounded-xl glass text-xs font-black uppercase tracking-[0.2em] hover:border-blue-500/50 transition-all">
             Scan Fleet
          </button>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="bg-zinc-50/50 dark:bg-white/[0.02] text-[10px] font-black text-zinc-500 uppercase tracking-[0.3em] border-b border-zinc-100 dark:border-white/5">
                <th className="px-10 py-5">Source Vector</th>
                <th className="px-10 py-5">Signal Type</th>
                <th className="px-10 py-5">Assigned Risk</th>
                <th className="px-10 py-5">Sentinel Action</th>
                <th className="px-10 py-5 text-right">Maturity</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-50 dark:divide-white/5 font-medium">
              {[
                { ip: '185.244.12.5', type: 'DNS_FLOOD', score: 85, action: 'QUARANTINE', time: '2m ago' },
                { ip: '45.10.2.198', type: 'SQLI_DETECT', score: 92, action: 'CLUSTER_BAN', time: '14m ago' },
                { ip: '102.33.1.4', type: 'SCAN_TCP', score: 65, action: 'RATE_LIMIT', time: '1h ago' },
              ].map((t, i) => (
                <tr key={i} className="group hover:bg-zinc-50 dark:hover:bg-white/[0.02] transition-colors cursor-pointer">
                  <td className="px-10 py-6">
                    <div className="flex items-center gap-4">
                       <div className="w-10 h-10 rounded-xl bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center">
                          <Globe className="w-5 h-5 text-zinc-400" />
                       </div>
                       <div>
                          <p className="font-bold dark:text-white text-zinc-950">{t.ip}</p>
                          <p className="text-[10px] text-zinc-500 font-mono">AS42031 • RO</p>
                       </div>
                    </div>
                  </td>
                  <td className="px-10 py-6 text-sm font-black text-zinc-600 dark:text-zinc-400 italic">
                    {t.type}
                  </td>
                  <td className="px-10 py-6">
                    <div className="w-32 h-2 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden">
                       <div className="h-full bg-rose-500" style={{ width: `${t.score}%` }} />
                    </div>
                  </td>
                  <td className="px-10 py-6">
                    <span className="px-3 py-1 rounded-lg bg-rose-500/10 text-rose-500 text-[10px] font-black uppercase tracking-widest border border-rose-500/20">
                      {t.action}
                    </span>
                  </td>
                  <td className="px-10 py-6 text-right text-xs text-zinc-400 font-bold">
                    {t.time}
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

export default Dashboard;

function cn(...inputs: any[]) {
  return inputs.filter(Boolean).join(' ');
}
