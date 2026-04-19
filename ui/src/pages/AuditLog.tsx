import React, { useState } from 'react';
import { 
  History, 
  Search, 
  Filter, 
  ShieldCheck, 
  ShieldAlert, 
  User, 
  Clock, 
  Database, 
  ArrowRight,
  Fingerprint,
  ChevronRight,
  ExternalLink,
  Download,
  AlertCircle
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const AuditLog: React.FC = () => {
  const [selectedEvent, setSelectedEvent] = useState<any>(null);

  const events = [
    { id: '1', timestamp: '2026-04-18 19:42:05', action: 'policy.apply', actor: 'ersin', resource: 'global-ingress', result: 'success', severity: 'info', hash: '8a3f...d21' },
    { id: '2', timestamp: '2026-04-18 19:40:12', action: 'ips.block', actor: 'sentinel', resource: '192.168.1.45', result: 'success', severity: 'high', hash: 'f2c1...a90' },
    { id: '3', timestamp: '2026-04-18 18:22:50', action: 'snapshot.create', actor: 'system', resource: 'pre-update', result: 'success', severity: 'info', hash: 'e3b9...c44' },
    { id: '4', timestamp: '2026-04-18 17:15:00', action: 'api.key_created', actor: 'admin', resource: 'ci-pipeline-key', result: 'success', severity: 'medium', hash: '99a1...ff2' },
    { id: '5', timestamp: '2026-04-18 16:05:12', action: 'policy.rollback', actor: 'ersin', resource: 'emergency-revert', result: 'failure', severity: 'high', hash: '12d3...e55' },
  ];

  return (
    <div className="space-y-8 animate-in fade-in duration-700">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div>
          <h1 className="text-4xl font-black tracking-tight dark:text-white text-zinc-950 flex items-center gap-4">
            <History className="w-10 h-10 text-blue-600" />
            Audit Forensics
          </h1>
          <p className="text-zinc-500 font-medium text-lg mt-1">Tamper-evident, cryptographically mühürlü event chain.</p>
        </div>
        <div className="flex gap-3 w-full md:w-auto">
           <button className="flex-1 md:flex-none inline-flex items-center justify-center gap-2 px-6 py-3 rounded-2xl bg-zinc-100 dark:bg-white/[0.05] text-zinc-600 dark:text-zinc-400 font-black text-sm border border-transparent hover:border-blue-500/50 transition-all group">
             <Fingerprint className="w-4 h-4 text-blue-500 group-hover:scale-110 transition-transform" />
             Verify Integrity
           </button>
           <button className="flex-1 md:flex-none inline-flex items-center justify-center gap-2 px-6 py-3 rounded-2xl bg-zinc-950 dark:bg-white text-white dark:text-zinc-950 font-black text-sm shadow-xl hover:opacity-90 transition-all">
             <Download className="w-4 h-4" />
             Export JSONL
           </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <div className="lg:col-span-3 bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] overflow-hidden shadow-sm">
          <div className="p-8 border-b border-zinc-100 dark:border-white/5 bg-zinc-50/50 dark:bg-white/[0.01] flex items-center gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-4 w-4 text-zinc-400" />
              <input 
                type="text" 
                placeholder="Query by ID, actor, or resource signature..." 
                className="w-full pl-11 pr-4 py-3 rounded-2xl bg-white dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 focus:ring-2 focus:ring-blue-500/50 outline-none text-sm"
              />
            </div>
            <button className="p-3 rounded-2xl bg-white dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 text-zinc-400 hover:text-blue-500 transition-all">
              <Filter className="w-5 h-5" />
            </button>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="bg-zinc-50/50 dark:bg-white/[0.02] text-[10px] font-black text-zinc-500 uppercase tracking-[0.3em] border-b border-zinc-100 dark:border-white/5">
                  <th className="px-10 py-5">Event Chain</th>
                  <th className="px-10 py-5">Security Action</th>
                  <th className="px-10 py-5">Resource Vector</th>
                  <th className="px-10 py-5">Mühür Hash</th>
                  <th className="px-10 py-5 text-right">Result</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-50 dark:divide-white/5">
                {events.map((event) => (
                  <tr 
                    key={event.id} 
                    onClick={() => setSelectedEvent(event)}
                    className={cn(
                      "group cursor-pointer transition-all duration-300",
                      selectedEvent?.id === event.id ? "bg-blue-600/10 dark:bg-blue-600/5" : "hover:bg-zinc-50 dark:hover:bg-white/[0.02]"
                    )}
                  >
                    <td className="px-10 py-6">
                      <div className="flex items-center gap-3">
                         <div className="w-2 h-2 rounded-full bg-blue-500/40 group-hover:scale-150 transition-transform" />
                         <div>
                            <p className="text-xs font-black dark:text-white text-zinc-950">{event.timestamp}</p>
                            <div className="flex items-center gap-2 mt-1">
                               <User className="w-3 h-3 text-zinc-400" />
                               <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">{event.actor}</span>
                            </div>
                         </div>
                      </div>
                    </td>
                    <td className="px-10 py-6">
                      <span className="text-sm font-black dark:text-zinc-300 text-zinc-700 italic">
                        {event.action}
                      </span>
                    </td>
                    <td className="px-10 py-6">
                      <div className="flex items-center gap-2 bg-zinc-100 dark:bg-white/5 px-3 py-1.5 rounded-xl border border-zinc-200/50 dark:border-white/5 w-fit">
                         <Database className="w-3 h-3 text-zinc-400" />
                         <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">{event.resource}</span>
                      </div>
                    </td>
                    <td className="px-10 py-6 font-mono text-[10px] text-zinc-400">
                      {event.hash}
                    </td>
                    <td className="px-10 py-6 text-right">
                      <span className={cn(
                        "px-3 py-1 rounded-lg text-[10px] font-black uppercase tracking-widest border",
                        event.result === 'success' 
                          ? 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20' 
                          : 'bg-rose-500/10 text-rose-600 border-rose-500/20'
                      )}>
                        {event.result}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Forensic Detail Panel */}
        <div className="lg:col-span-1 space-y-6">
           <AnimatePresence mode="wait">
             {selectedEvent ? (
               <motion.div 
                 key="selected"
                 initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: 20 }}
                 className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] p-8 shadow-sm flex flex-col h-full sticky top-8"
               >
                  <div className="flex items-center justify-between mb-8">
                     <h3 className="text-lg font-black dark:text-white text-zinc-950">Forensic View</h3>
                     <button onClick={() => setSelectedEvent(null)} className="p-2 hover:bg-zinc-100 dark:hover:bg-white/5 rounded-xl transition-all">
                        <X className="w-4 h-4 text-zinc-400" />
                     </button>
                  </div>

                  <div className="space-y-8 flex-1 overflow-y-auto pr-2">
                     <div className="p-6 rounded-[2rem] bg-zinc-50 dark:bg-zinc-950 border border-zinc-100 dark:border-zinc-800">
                        <div className="flex items-center gap-3 mb-4 text-blue-600">
                           <Fingerprint className="w-5 h-5" />
                           <span className="text-[10px] font-black uppercase tracking-widest">Evidence ID</span>
                        </div>
                        <p className="font-mono text-xs dark:text-zinc-300 break-all">{selectedEvent.id}-{selectedEvent.hash.replace('...', '')}</p>
                     </div>

                     <div className="space-y-4">
                        <h4 className="text-[10px] font-black text-zinc-500 uppercase tracking-widest ml-1">Contextual Change</h4>
                        <div className="p-5 rounded-2xl bg-zinc-900 border border-zinc-800 font-mono text-[11px] leading-relaxed relative group">
                           <div className="text-rose-400 line-through">- protocol: tcp</div>
                           <div className="text-emerald-400 font-bold">+ protocol: any</div>
                           <div className="text-zinc-500">  destPorts: [80, 443]</div>
                           <button className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity text-zinc-400 hover:text-white">
                              <Code2 className="w-4 h-4" />
                           </button>
                        </div>
                     </div>

                     <div className="p-6 rounded-[2.5rem] bg-blue-600/5 border border-blue-500/10">
                        <div className="flex items-center gap-3 mb-3">
                           <ShieldCheck className="w-5 h-5 text-emerald-500" />
                           <span className="text-xs font-black dark:text-white uppercase tracking-tighter text-zinc-950 italic">Verified State</span>
                        </div>
                        <p className="text-xs text-zinc-500 font-bold leading-relaxed">This event is linked to the previous 1,248 events via a secure hash-chain.</p>
                     </div>
                  </div>

                  <div className="mt-8 pt-8 border-t border-zinc-100 dark:border-white/5 space-y-4">
                     <button className="w-full py-4 rounded-2xl bg-blue-600 hover:bg-blue-500 text-white font-black shadow-xl shadow-blue-600/30 flex items-center justify-center gap-3 transition-all text-sm">
                        Inspect Raw Payload <ArrowRight className="w-4 h-4" />
                     </button>
                  </div>
               </motion.div>
             ) : (
               <motion.div 
                 key="empty"
                 initial={{ opacity: 0 }} animate={{ opacity: 1 }}
                 className="h-full border-2 border-dashed border-zinc-200 dark:border-white/5 rounded-[3rem] flex flex-col items-center justify-center p-12 text-center"
               >
                  <div className="w-20 h-20 rounded-full bg-zinc-50 dark:bg-white/[0.02] flex items-center justify-center mb-6">
                     <Search className="w-10 h-10 text-zinc-200 dark:text-zinc-800" />
                  </div>
                  <h4 className="text-xl font-black dark:text-zinc-400 text-zinc-300">No Event Selected</h4>
                  <p className="text-sm text-zinc-500 font-bold mt-2">Select an audit entry from the forensic chain to inspect signatures.</p>
               </motion.div>
             )}
           </AnimatePresence>
        </div>
      </div>
    </div>
  );
};

export default AuditLog;

function cn(...inputs: any[]) {
  return inputs.filter(Boolean).join(' ');
}
