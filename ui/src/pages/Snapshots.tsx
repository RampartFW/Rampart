import React, { useState } from 'react';
import { 
  Camera, 
  RotateCcw, 
  History, 
  Zap, 
  Shield, 
  ArrowLeft, 
  ArrowRight,
  Database,
  Search,
  Clock,
  ChevronRight,
  Plus,
  Trash2,
  FileText
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const Snapshots: React.FC = () => {
  const [selectedSnapshot, setSelectedSnapshot] = useState<any>(null);

  const snapshots = [
    { id: 'snap-8812', timestamp: '2026-04-18 19:42:05', trigger: 'policy.apply', description: 'Automatic snapshot before global-ingress update', rules: 156, size: '42KB' },
    { id: 'snap-8790', timestamp: '2026-04-18 10:15:22', trigger: 'manual', description: 'Pre-maintenance backup (Ersin)', rules: 154, size: '41KB' },
    { id: 'snap-8744', timestamp: '2026-04-17 22:30:10', trigger: 'system', description: 'Daily scheduled state backup', rules: 154, size: '41KB' },
    { id: 'snap-8611', timestamp: '2026-04-16 14:05:00', trigger: 'ips.ban', description: 'Auto-snapshot after massive IPS quarantine', rules: 158, size: '43KB' },
  ];

  return (
    <div className="space-y-8 animate-in fade-in duration-700">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div>
          <h1 className="text-4xl font-black tracking-tight dark:text-white text-zinc-950 flex items-center gap-4">
            <Camera className="w-10 h-10 text-blue-600" />
            Infrastructure Snapshots
          </h1>
          <p className="text-zinc-500 font-medium text-lg mt-1">Immutable state captures for instant disaster recovery.</p>
        </div>
        <div className="flex gap-3 w-full md:w-auto">
           <button className="flex-1 md:flex-none inline-flex items-center justify-center gap-2 px-6 py-3 rounded-2xl bg-blue-600 text-white font-black text-sm shadow-xl shadow-blue-600/30 hover:bg-blue-500 transition-all">
             <Plus className="w-4 h-4" />
             Capture State
           </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
           {snapshots.map((snap, i) => (
             <motion.div 
               key={snap.id}
               initial={{ opacity: 0, x: -20 }}
               animate={{ opacity: 1, x: 0 }}
               transition={{ delay: i * 0.1 }}
               onClick={() => setSelectedSnapshot(snap)}
               className={cn(
                 "p-8 rounded-[2.5rem] border transition-all duration-500 cursor-pointer group relative overflow-hidden",
                 selectedSnapshot?.id === snap.id 
                   ? "bg-blue-600 border-blue-500 shadow-2xl shadow-blue-600/20 translate-x-4" 
                   : "bg-white dark:bg-zinc-900/50 border-zinc-200 dark:border-white/5 hover:border-blue-500/40"
               )}
             >
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 relative z-10">
                   <div className="flex items-center gap-6">
                      <div className={cn(
                        "w-16 h-16 rounded-2xl flex items-center justify-center shadow-lg transition-transform duration-500 group-hover:rotate-12",
                        selectedSnapshot?.id === snap.id ? "bg-white/20 text-white" : "bg-zinc-100 dark:bg-zinc-800 text-zinc-400"
                      )}>
                        <Database className="w-7 h-7" />
                      </div>
                      <div>
                         <div className="flex items-center gap-3">
                            <h3 className={cn(
                              "text-xl font-black tracking-tight",
                              selectedSnapshot?.id === snap.id ? "text-white" : "dark:text-white text-zinc-950"
                            )}>{snap.id}</h3>
                            <span className={cn(
                              "px-2 py-0.5 rounded-lg text-[10px] font-black uppercase tracking-widest border",
                              selectedSnapshot?.id === snap.id 
                                ? "bg-white/10 border-white/20 text-white" 
                                : "bg-zinc-100 dark:bg-white/5 border-zinc-200 dark:border-white/10 text-zinc-500"
                            )}>
                              {snap.trigger}
                            </span>
                         </div>
                         <p className={cn(
                           "text-sm font-bold mt-1",
                           selectedSnapshot?.id === snap.id ? "text-blue-100" : "text-zinc-500"
                         )}>{snap.description}</p>
                      </div>
                   </div>
                   
                   <div className="flex items-center gap-8 text-right shrink-0">
                      <div className="hidden sm:block">
                         <p className={cn(
                           "text-xs font-black uppercase tracking-widest",
                           selectedSnapshot?.id === snap.id ? "text-blue-100" : "text-zinc-400"
                         )}>Rules</p>
                         <p className={cn(
                           "text-lg font-black",
                           selectedSnapshot?.id === snap.id ? "text-white" : "dark:text-zinc-200"
                         )}>{snap.rules}</p>
                      </div>
                      <div className="hidden sm:block">
                         <p className={cn(
                           "text-xs font-black uppercase tracking-widest",
                           selectedSnapshot?.id === snap.id ? "text-blue-100" : "text-zinc-400"
                         )}>Captured</p>
                         <p className={cn(
                           "text-xs font-bold",
                           selectedSnapshot?.id === snap.id ? "text-white" : "dark:text-zinc-400"
                         )}>{snap.timestamp.split(' ')[1]}</p>
                      </div>
                      <ChevronRight className={cn(
                        "w-5 h-5 transition-transform duration-500",
                        selectedSnapshot?.id === snap.id ? "text-white rotate-90" : "text-zinc-300 dark:text-zinc-700 group-hover:translate-x-1"
                      )} />
                   </div>
                </div>
                {selectedSnapshot?.id === snap.id && (
                  <div className="absolute top-0 right-0 w-64 h-full bg-gradient-to-l from-white/10 to-transparent pointer-events-none" />
                )}
             </motion.div>
           ))}
        </div>

        {/* Snapshot Utility Panel */}
        <div className="lg:col-span-1">
           <AnimatePresence mode="wait">
              {selectedSnapshot ? (
                <motion.div 
                  key="utility"
                  initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.95 }}
                  className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] p-10 shadow-sm flex flex-col sticky top-8"
                >
                   <h3 className="text-2xl font-black dark:text-white text-zinc-950 mb-8 tracking-tighter">Time Machine</h3>
                   
                   <div className="space-y-8 flex-1">
                      <div className="p-8 rounded-[2.5rem] bg-zinc-950 border border-white/10 text-center group">
                         <RotateCcw className="w-12 h-12 text-rose-500 mx-auto mb-4 group-hover:rotate-[-45deg] transition-transform duration-500" />
                         <h4 className="text-white font-black text-xl mb-2">Rollback Point</h4>
                         <p className="text-xs text-zinc-500 font-bold leading-relaxed px-4 text-zinc-400">Reverting to this state will atomicly overwrite the current kernel tables.</p>
                      </div>

                      <div className="space-y-4">
                         <div className="flex items-center justify-between px-2">
                            <span className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Snapshot Stats</span>
                            <span className="text-[10px] font-mono text-zinc-400 italic">verified checksum</span>
                         </div>
                         <div className="grid grid-cols-2 gap-4">
                            <div className="p-4 rounded-2xl bg-zinc-50 dark:bg-white/5 border border-zinc-100 dark:border-white/10">
                               <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest mb-1">Disk Size</p>
                               <p className="text-sm font-black dark:text-white">{selectedSnapshot.size}</p>
                            </div>
                            <div className="p-4 rounded-2xl bg-zinc-50 dark:bg-white/5 border border-zinc-100 dark:border-white/10">
                               <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest mb-1">Format</p>
                               <p className="text-sm font-black dark:text-white italic underline">gob/zstd</p>
                            </div>
                         </div>
                      </div>
                   </div>

                   <div className="mt-12 space-y-4">
                      <button className="w-full py-5 rounded-[1.5rem] bg-rose-600 hover:bg-rose-500 text-white font-black shadow-2xl shadow-rose-600/30 flex items-center justify-center gap-3 transition-all">
                         <RotateCcw className="w-5 h-5" />
                         Apply Rollback
                      </button>
                      <div className="grid grid-cols-2 gap-4">
                        <button className="py-4 rounded-2xl glass text-zinc-500 hover:text-blue-600 font-bold text-xs flex items-center justify-center gap-2">
                           <FileText className="w-4 h-4" />
                           Diff YAML
                        </button>
                        <button className="py-4 rounded-2xl glass text-zinc-500 hover:text-rose-600 font-bold text-xs flex items-center justify-center gap-2">
                           <Trash2 className="w-4 h-4" />
                           Purge
                        </button>
                      </div>
                   </div>
                </motion.div>
              ) : (
                <div className="h-[400px] border-2 border-dashed border-zinc-200 dark:border-white/5 rounded-[4rem] flex flex-col items-center justify-center p-12 text-center">
                   <div className="w-20 h-20 rounded-full bg-zinc-50 dark:bg-white/[0.02] flex items-center justify-center mb-6">
                      <History className="w-10 h-10 text-zinc-200 dark:text-zinc-800" />
                   </div>
                   <h4 className="text-xl font-black dark:text-zinc-400 text-zinc-300">Target a State</h4>
                   <p className="text-sm text-zinc-500 font-bold mt-2 leading-relaxed px-6">Choose a snapshot from the timeline to initiate a differential analysis or rollback.</p>
                </div>
              )}
           </AnimatePresence>
        </div>
      </div>
    </div>
  );
};

export default Snapshots;

function cn(...inputs: any[]) {
  return inputs.filter(Boolean).join(' ');
}
