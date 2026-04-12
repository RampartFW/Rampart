import React from 'react';
import { Zap, CheckCircle2 } from 'lucide-react';

const Simulator: React.FC = () => {
  return (
    <div className="space-y-6 max-w-5xl mx-auto">
      <div>
        <h2 className="text-2xl font-bold">Packet Simulator</h2>
        <p className="text-sm text-gray-500">Test if a packet would be allowed or denied by current rules</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card p-6">
          <h3 className="text-lg font-medium mb-4">Packet Configuration</h3>
          <form className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">Source IP</label>
                <input type="text" className="input" placeholder="e.g. 10.0.1.50" defaultValue="10.0.1.50" />
              </div>
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">Dest IP</label>
                <input type="text" className="input" placeholder="e.g. 192.168.1.10" defaultValue="192.168.1.10" />
              </div>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">Protocol</label>
                <select className="input">
                  <option>TCP</option>
                  <option>UDP</option>
                  <option>ICMP</option>
                  <option>ICMPv6</option>
                </select>
              </div>
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">Dest Port</label>
                <input type="number" className="input" placeholder="e.g. 80" defaultValue="22" />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">Direction</label>
                <select className="input">
                  <option>Inbound</option>
                  <option>Outbound</option>
                  <option>Forward</option>
                </select>
              </div>
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">Interface</label>
                <input type="text" className="input" placeholder="e.g. eth0" defaultValue="eth0" />
              </div>
            </div>

            <button type="button" className="btn btn-primary w-full py-3 mt-4">
              <Zap className="h-4 w-4 mr-2" />
              Run Simulation
            </button>
          </form>
        </div>

        <div className="space-y-6">
          <div className="card p-6 border-l-4 border-green-500 bg-green-50/50 dark:bg-green-900/10">
            <div className="flex items-center gap-4">
              <CheckCircle2 className="h-10 w-10 text-green-500" />
              <div>
                <h3 className="text-2xl font-bold text-green-700 dark:text-green-400">ALLOWED</h3>
                <p className="text-sm text-green-600 dark:text-green-500/80">Packet matches rule "allow-ssh-bastion"</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <h3 className="text-sm font-bold mb-4 uppercase tracking-wider text-gray-500">Evaluation Trace</h3>
            <div className="space-y-4">
              <div className="relative pl-6 pb-4 border-l border-gray-200 dark:border-gray-700 last:border-0 last:pb-0">
                <div className="absolute left-[-5px] top-1 h-2 w-2 rounded-full bg-gray-300 dark:bg-gray-600"></div>
                <div className="text-xs font-medium text-gray-500">Rule [P005] system-default-drop</div>
                <div className="text-sm text-gray-400 italic">No match: direction mismatch</div>
              </div>
              <div className="relative pl-6 pb-4 border-l border-gray-200 dark:border-gray-700 last:border-0 last:pb-0">
                <div className="absolute left-[-5px] top-1 h-2 w-2 rounded-full bg-green-500"></div>
                <div className="text-xs font-bold text-green-600">Rule [P010] allow-ssh-bastion</div>
                <div className="text-sm font-medium">MATCH: protocol=tcp, dport=22, source=10.0.1.0/24</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Simulator;
