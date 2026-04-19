import React from 'react';
import { 
  Settings as SettingsIcon, 
  Key, 
  Shield, 
  Database, 
  RefreshCw, 
  Bell, 
  Monitor, 
  Terminal,
  Save,
  CheckCircle2,
  Lock,
  Cpu
} from 'lucide-react';

const Settings: React.FC = () => {
  return (
    <div className="max-w-5xl space-y-10 animate-in fade-in duration-700">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div>
          <h1 className="text-4xl font-black tracking-tight dark:text-white text-zinc-950 flex items-center gap-4">
            <SettingsIcon className="w-10 h-10 text-blue-600" />
            System Control
          </h1>
          <p className="text-zinc-500 font-medium text-lg mt-1">Configure global parameters and security credentials.</p>
        </div>
        <button className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-2xl bg-blue-600 text-white font-black text-sm shadow-xl shadow-blue-600/30 hover:bg-blue-500 transition-all active:scale-95">
          <Save className="w-4 h-4" />
          Update Parameters
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-10">
         {/* Navigation Sidebar */}
         <div className="lg:col-span-1 space-y-2">
            {[
              { id: 'auth', name: 'Authentication', icon: Key, active: true },
              { id: 'nodes', name: 'Node Defaults', icon: Cpu, active: false },
              { id: 'audit', name: 'Audit & Retention', icon: Database, active: false },
              { id: 'alerts', name: 'Alert Channels', icon: Bell, active: false },
              { id: 'display', name: 'Interface Prefs', icon: Monitor, active: false },
            ].map(item => (
              <button 
                key={item.id}
                className={cn(
                  "w-full flex items-center gap-4 px-6 py-4 rounded-2xl font-bold text-sm transition-all",
                  item.active 
                    ? "bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-white/10 shadow-sm dark:text-white" 
                    : "text-zinc-500 hover:bg-zinc-100 dark:hover:bg-white/5"
                )}
              >
                <item.icon className={cn("w-5 h-5", item.active ? "text-blue-500" : "text-zinc-400")} />
                {item.name}
              </button>
            ))}
         </div>

         {/* Content Area */}
         <div className="lg:col-span-2 space-y-10">
            <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] p-10 shadow-sm">
               <h3 className="text-xl font-black dark:text-white text-zinc-950 mb-8 flex items-center gap-3">
                  <Key className="w-5 h-5 text-blue-500" />
                  API Access Keys
               </h3>
               
               <div className="space-y-6">
                  {[
                    { name: 'Default Admin', prefix: 'rmp_live_***', created: '2 days ago', status: 'active' },
                    { name: 'CI/CD Pipeline', prefix: 'rmp_test_***', created: '14 days ago', status: 'active' },
                  ].map((key, i) => (
                    <div key={i} className="flex items-center justify-between p-6 rounded-2xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-100 dark:border-zinc-800 group">
                       <div className="flex items-center gap-4">
                          <div className="w-10 h-10 rounded-xl bg-white dark:bg-zinc-900 flex items-center justify-center border border-zinc-200 dark:border-zinc-800">
                             <Shield className="w-5 h-5 text-zinc-400" />
                          </div>
                          <div>
                             <p className="font-bold dark:text-white text-zinc-900">{key.name}</p>
                             <p className="text-xs font-mono text-zinc-500">{key.prefix}</p>
                          </div>
                       </div>
                       <div className="flex items-center gap-6">
                          <div className="text-right hidden sm:block">
                             <p className="text-[10px] font-black uppercase text-zinc-400">Created</p>
                             <p className="text-xs font-bold text-zinc-500">{key.created}</p>
                          </div>
                          <span className="px-2 py-0.5 rounded-md bg-emerald-500/10 text-emerald-500 text-[10px] font-black uppercase border border-emerald-500/20">
                             {key.status}
                          </span>
                       </div>
                    </div>
                  ))}
                  <button className="w-full py-4 rounded-2xl border-2 border-dashed border-zinc-200 dark:border-zinc-800 text-zinc-400 hover:border-blue-500/50 hover:text-blue-500 font-bold text-sm transition-all flex items-center justify-center gap-2">
                     <Plus className="w-4 h-4" />
                     Generate New Access Key
                  </button>
               </div>
            </div>

            <div className="bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-white/5 rounded-[3rem] p-10 shadow-sm">
               <h3 className="text-xl font-black dark:text-white text-zinc-950 mb-8 flex items-center gap-3">
                  <Database className="w-5 h-5 text-indigo-500" />
                  Persistence Settings
               </h3>
               
               <div className="space-y-8">
                  <div className="space-y-4">
                     <label className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em] ml-2">Snapshot Strategy</label>
                     <div className="grid grid-cols-2 gap-4">
                        <button className="p-4 rounded-2xl bg-blue-600/5 border border-blue-500/20 text-left group">
                           <div className="flex justify-between items-center mb-1">
                              <span className="font-bold text-sm text-blue-600">On Change</span>
                              <CheckCircle2 className="w-4 h-4 text-blue-500" />
                           </div>
                           <p className="text-[10px] text-zinc-500 font-medium">Capture state before every apply</p>
                        </button>
                        <button className="p-4 rounded-2xl bg-zinc-50 dark:bg-white/5 border border-zinc-100 dark:border-white/10 text-left opacity-60">
                           <div className="flex justify-between items-center mb-1">
                              <span className="font-bold text-sm text-zinc-500">Periodic</span>
                           </div>
                           <p className="text-[10px] text-zinc-500 font-medium">Capture every 24 hours</p>
                        </button>
                     </div>
                  </div>

                  <div className="space-y-2">
                     <div className="flex justify-between items-center px-2">
                        <label className="text-[10px] font-black text-zinc-500 uppercase tracking-[0.2em]">Retention Limit</label>
                        <span className="text-sm font-black text-blue-500">30 Snapshots</span>
                     </div>
                     <input type="range" className="w-full h-1.5 bg-zinc-100 dark:bg-zinc-800 rounded-full appearance-none cursor-pointer accent-blue-600" />
                  </div>
               </div>
            </div>
         </div>
      </div>
    </div>
  );
};

export default Settings;

function cn(...inputs: any[]) {
  return inputs.filter(Boolean).join(' ');
}

function Plus({ className }: any) {
  return <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>;
}
