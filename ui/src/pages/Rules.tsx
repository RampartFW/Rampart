import React, { useState } from 'react';
import { 
  Search, 
  Plus, 
  BarChart2, 
  Filter, 
  MoreVertical, 
  Shield, 
  ArrowRight,
  Zap,
  Globe,
  Database
} from 'lucide-react';

const Rules: React.FC = () => {
  const [searchTerm, setSearch] = useState('');
  
  const rules = [
    { id: '1', name: 'allow-ssh-bastion', priority: 10, protocol: 'tcp', dport: '22', source: '10.0.1.0/24', action: 'accept', pkts: '1.2k', bytes: '840KB', backend: 'nftables' },
    { id: '2', name: 'allow-http', priority: 500, protocol: 'tcp', dport: '80', source: '0.0.0.0/0', action: 'accept', pkts: '124k', bytes: '12MB', backend: 'eBPF' },
    { id: '3', name: 'allow-https', priority: 500, protocol: 'tcp', dport: '443', source: '0.0.0.0/0', action: 'accept', pkts: '450k', bytes: '84MB', backend: 'eBPF' },
    { id: '4', name: 'deny-all-ssh', priority: 10, protocol: 'tcp', dport: '22', source: '0.0.0.0/0', action: 'drop', pkts: '4.5k', bytes: '2.4MB', backend: 'nftables' },
    { id: '5', name: 'aws-db-access', priority: 300, protocol: 'tcp', dport: '5432', source: '10.0.5.0/24', action: 'accept', pkts: '0', bytes: '0', backend: 'AWS SG' },
  ];

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight dark:text-white">Active Rules</h1>
          <p className="text-zinc-500 mt-1">Orchestrated view of every filter across your distributed infrastructure.</p>
        </div>
        <div className="flex gap-3 w-full md:w-auto">
           <button className="flex-1 md:flex-none inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-zinc-100 dark:bg-white/[0.05] text-zinc-600 dark:text-zinc-400 font-bold text-sm border border-transparent hover:border-zinc-300 dark:hover:border-zinc-700 transition-all">
             <Filter className="w-4 h-4" />
             Filters
           </button>
           <button className="flex-1 md:flex-none inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-blue-600 text-white font-bold text-sm shadow-lg shadow-blue-500/25 hover:bg-blue-500 transition-all">
             <Plus className="w-4 h-4" />
             Quick Add
           </button>
        </div>
      </div>

      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl overflow-hidden shadow-sm">
        <div className="p-6 border-b border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-white/[0.01]">
          <div className="relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-4 w-4 text-zinc-400" />
            <input 
              type="text" 
              placeholder="Search by name, CIDR, port or backend..." 
              value={searchTerm}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-11 pr-4 py-3 rounded-2xl bg-white dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all text-sm"
            />
          </div>
        </div>
        
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-zinc-50/50 dark:bg-white/[0.02] text-[10px] font-bold text-zinc-500 uppercase tracking-widest border-b border-zinc-100 dark:border-zinc-800">
                <th className="px-8 py-4">Priority</th>
                <th className="px-8 py-4">Rule Identity</th>
                <th className="px-8 py-4">Orchestration</th>
                <th className="px-8 py-4">Conditions</th>
                <th className="px-8 py-4">Verdict</th>
                <th className="px-8 py-4">Live Traffic</th>
                <th className="px-8 py-4 text-right"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-50 dark:divide-zinc-800/50">
              {rules.map((rule) => (
                <tr key={rule.id} className="group hover:bg-zinc-50 dark:hover:bg-white/[0.01] transition-colors">
                  <td className="px-8 py-6 font-mono text-xs text-zinc-400">
                    {rule.priority.toString().padStart(3, '0')}
                  </td>
                  <td className="px-8 py-6">
                    <div className="font-bold text-zinc-900 dark:text-zinc-100">{rule.name}</div>
                    <div className="text-[10px] text-zinc-500 font-mono mt-0.5">ID: {rule.id.padStart(4, '0')}</div>
                  </td>
                  <td className="px-8 py-6">
                    <div className="flex items-center gap-2">
                       <div className={cn(
                         "w-1.5 h-1.5 rounded-full animate-pulse",
                         rule.backend === 'eBPF' ? 'bg-blue-500' : rule.backend === 'AWS SG' ? 'bg-orange-500' : 'bg-emerald-500'
                       )} />
                       <span className="text-xs font-bold text-zinc-600 dark:text-zinc-400">{rule.backend}</span>
                    </div>
                  </td>
                  <td className="px-8 py-6">
                    <div className="flex flex-wrap gap-1.5">
                      <span className="px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-white/[0.05] text-[10px] font-extrabold uppercase tracking-tight text-zinc-500">{rule.protocol}</span>
                      <span className="px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-white/[0.05] text-[10px] font-mono text-zinc-500">:{rule.dport}</span>
                      <span className="inline-flex items-center px-2 py-0.5 rounded-md bg-blue-500/5 text-blue-600 dark:text-blue-400 text-[10px] font-mono border border-blue-500/10">
                        <ArrowRight className="w-2.5 h-2.5 mr-1" />
                        {rule.source}
                      </span>
                    </div>
                  </td>
                  <td className="px-8 py-6">
                    <span className={cn(
                      "px-3 py-1 text-[10px] font-extrabold rounded-lg uppercase tracking-widest border",
                      rule.action === 'accept' 
                        ? 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20' 
                        : 'bg-rose-500/10 text-rose-600 border-rose-500/20'
                    )}>
                      {rule.action}
                    </span>
                  </td>
                  <td className="px-8 py-6">
                    <div className="flex items-center gap-4">
                      <div className="flex flex-col">
                        <span className="text-xs font-bold dark:text-zinc-200">{rule.pkts} pkts</span>
                        <span className="text-[10px] text-zinc-500">{rule.bytes} transferred</span>
                      </div>
                      <div className="h-6 w-[1px] bg-zinc-100 dark:bg-zinc-800" />
                      <BarChart2 className="h-4 w-4 text-zinc-300 dark:text-zinc-700" />
                    </div>
                  </td>
                  <td className="px-8 py-6 text-right">
                    <button className="p-2 rounded-lg text-zinc-400 hover:bg-zinc-100 dark:hover:bg-white/[0.05] transition-all">
                      <MoreVertical className="h-4 w-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
      
      <div className="grid md:grid-cols-3 gap-6">
         <div className="p-6 rounded-3xl bg-blue-600/5 border border-blue-500/10 flex items-center gap-4">
            <div className="w-12 h-12 rounded-2xl bg-blue-600 flex items-center justify-center shadow-lg shadow-blue-500/20">
               <Zap className="w-6 h-6 text-white" />
            </div>
            <div>
               <p className="text-[10px] font-bold text-blue-500 uppercase tracking-widest mb-0.5">XDP Fast-Path</p>
               <h4 className="text-sm font-bold dark:text-white text-zinc-900">85% of traffic offloaded</h4>
            </div>
         </div>
         <div className="p-6 rounded-3xl bg-orange-600/5 border border-orange-500/10 flex items-center gap-4">
            <div className="w-12 h-12 rounded-2xl bg-orange-600 flex items-center justify-center shadow-lg shadow-orange-500/20">
               <Cloud className="w-6 h-6 text-white" />
            </div>
            <div>
               <p className="text-[10px] font-bold text-orange-500 uppercase tracking-widest mb-0.5">Multi-Cloud</p>
               <h4 className="text-sm font-bold dark:text-white text-zinc-900">SG Quorum Healthy</h4>
            </div>
         </div>
         <div className="p-6 rounded-3xl bg-indigo-600/5 border border-indigo-500/10 flex items-center gap-4">
            <div className="w-12 h-12 rounded-2xl bg-indigo-600 flex items-center justify-center shadow-lg shadow-indigo-500/20">
               <Shield className="w-6 h-6 text-white" />
            </div>
            <div>
               <p className="text-[10px] font-bold text-indigo-500 uppercase tracking-widest mb-0.5">Integrity</p>
               <h4 className="text-sm font-bold dark:text-white text-zinc-900">Mühürlü (Hash-chain) OK</h4>
            </div>
         </div>
      </div>
    </div>
  );
};

export default Rules;

function cn(...inputs: any[]) {
  return inputs.filter(Boolean).join(' ');
}
