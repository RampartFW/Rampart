import React, { useState } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import { yaml } from '@codemirror/lang-yaml';
import { 
  Save, 
  Play, 
  AlertTriangle, 
  History, 
  CheckCircle2, 
  Info,
  Shield,
  Layers,
  Clock
} from 'lucide-react';
import { motion } from 'framer-motion';

const Policies: React.FC = () => {
  const [code, setCode] = useState(`# rampart-policy.yaml
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: production-web-tier
  description: "Firewall rules for production web servers"

defaults:
  direction: inbound
  action: drop
  ipVersion: both
  states: [established, related]

policies:
  - name: ssh-access
    priority: 10
    rules:
      - name: allow-ssh-bastion
        match:
          protocol: tcp
          destPorts: [22]
          sourceCIDRs: ["10.0.1.0/24"]
        action: accept
`);

  return (
    <div className="h-full flex flex-col gap-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight dark:text-white">Policy Architect</h1>
          <p className="text-zinc-500 mt-1">Declarative network intent with real-time conflict detection.</p>
        </div>
        <div className="flex gap-3 w-full md:w-auto">
          <button className="flex-1 md:flex-none inline-flex items-center justify-center gap-2 px-6 py-2.5 rounded-xl bg-zinc-100 dark:bg-white/[0.05] text-zinc-600 dark:text-zinc-400 font-bold text-sm border border-transparent hover:border-zinc-300 dark:hover:border-zinc-700 transition-all">
            <Play className="h-4 w-4" />
            Dry Run
          </button>
          <button className="flex-1 md:flex-none inline-flex items-center justify-center gap-2 px-6 py-2.5 rounded-xl bg-blue-600 text-white font-bold text-sm shadow-lg shadow-blue-500/25 hover:bg-blue-500 transition-all">
            <Save className="h-4 w-4" />
            Deploy Policy
          </button>
        </div>
      </div>

      <div className="flex-1 grid grid-cols-1 lg:grid-cols-4 gap-8 min-h-0">
        <div className="lg:col-span-3 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl overflow-hidden shadow-sm flex flex-col min-h-[500px]">
          <div className="px-6 py-4 border-b border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-white/[0.01] flex items-center justify-between">
             <div className="flex items-center gap-2">
                <div className="flex gap-1.5 mr-4">
                  <div className="w-2.5 h-2.5 rounded-full bg-rose-500/40" />
                  <div className="w-2.5 h-2.5 rounded-full bg-amber-500/40" />
                  <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/40" />
                </div>
                <span className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest font-mono">active-policy.yaml</span>
             </div>
             <div className="flex items-center gap-4">
                <span className="text-[10px] font-bold text-emerald-500 uppercase tracking-widest flex items-center gap-1.5">
                  <CheckCircle2 className="w-3 h-3" />
                  Syntax Valid
                </span>
             </div>
          </div>
          <div className="flex-1 overflow-auto">
            <CodeMirror
              value={code}
              height="100%"
              theme="dark"
              extensions={[yaml()]}
              onChange={(value) => setCode(value)}
              className="text-sm h-full"
              basicSetup={{
                lineNumbers: true,
                foldGutter: true,
                highlightActiveLine: true,
              }}
            />
          </div>
        </div>
        
        <div className="space-y-6 flex flex-col">
          <div className="bg-amber-500/5 border border-amber-500/20 p-6 rounded-[2rem] relative overflow-hidden group">
            <div className="flex items-start gap-4 relative z-10">
              <AlertTriangle className="h-6 w-6 text-amber-500 shrink-0" />
              <div>
                <h4 className="text-sm font-bold text-amber-600 dark:text-amber-400">Architectural Warning</h4>
                <p className="mt-2 text-xs text-amber-700/80 dark:text-amber-500/70 leading-relaxed">
                  Rule <span className="font-mono text-amber-600 dark:text-amber-300">"deny-ssh-all"</span> is partially shadowed by <span className="font-mono text-amber-600 dark:text-amber-300">"allow-ssh-bastion"</span>. 
                  The traffic from <span className="font-mono underline">10.0.1.0/24</span> will never reach the drop rule.
                </p>
              </div>
            </div>
            <div className="absolute -bottom-4 -right-4 w-16 h-16 bg-amber-500/10 blur-2xl rounded-full group-hover:scale-150 transition-transform duration-700" />
          </div>
          
          <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 p-8 rounded-[2rem] shadow-sm flex-1">
            <h3 className="text-xs font-bold mb-6 uppercase tracking-widest text-zinc-400 flex items-center gap-2">
              <Layers className="w-4 h-4" />
              Execution Plan
            </h3>
            <div className="space-y-6">
              <div className="flex items-center justify-between group cursor-help">
                <div className="flex items-center gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]" />
                  <span className="text-xs font-bold text-zinc-600 dark:text-zinc-300">Rules to add</span>
                </div>
                <span className="text-sm font-mono font-bold text-emerald-500">+2</span>
              </div>
              <div className="flex items-center justify-between group cursor-help">
                <div className="flex items-center gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.5)]" />
                  <span className="text-xs font-bold text-zinc-600 dark:text-zinc-300">Rules to remove</span>
                </div>
                <span className="text-sm font-mono font-bold text-rose-500">-1</span>
              </div>
              <div className="flex items-center justify-between group cursor-help">
                <div className="flex items-center gap-3">
                  <div className="w-1.5 h-1.5 rounded-full bg-blue-500" />
                  <span className="text-xs font-bold text-zinc-600 dark:text-zinc-300">Modified</span>
                </div>
                <span className="text-sm font-mono font-bold text-blue-500">0</span>
              </div>
            </div>

            <div className="mt-10 pt-10 border-t border-zinc-100 dark:border-zinc-800">
               <h4 className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest mb-4 flex items-center gap-2">
                 <History className="w-3 h-3" />
                 Recent History
               </h4>
               <div className="space-y-4">
                  {[
                    { user: 'ersin', time: '14m ago', status: 'success' },
                    { user: 'system', time: '2h ago', status: 'success' },
                  ].map((h, i) => (
                    <div key={i} className="flex items-center justify-between">
                       <span className="text-xs font-medium text-zinc-500 dark:text-zinc-400">{h.user}</span>
                       <span className="text-[10px] font-mono text-zinc-400">{h.time}</span>
                    </div>
                  ))}
               </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Policies;
