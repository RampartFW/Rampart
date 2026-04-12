import React, { useState } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import { yaml } from '@codemirror/lang-yaml';
import { Save, Play, AlertTriangle } from 'lucide-react';

const Policies: React.FC = () => {
  const [code, setCode] = useState(`# rampart-policy.yaml
apiVersion: rampart.dev/v1
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
    <div className="h-full flex flex-col gap-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold">Policy Editor</h2>
          <p className="text-sm text-gray-500">Edit and apply firewall policies using YAML</p>
        </div>
        <div className="flex gap-3">
          <button className="btn btn-secondary">
            <Play className="h-4 w-4 mr-2" />
            Plan
          </button>
          <button className="btn btn-primary">
            <Save className="h-4 w-4 mr-2" />
            Apply Policy
          </button>
        </div>
      </div>

      <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-6 overflow-hidden">
        <div className="lg:col-span-2 card flex flex-col overflow-hidden">
          <CodeMirror
            value={code}
            height="100%"
            theme="dark"
            extensions={[yaml()]}
            onChange={(value) => setCode(value)}
            className="flex-1 text-sm overflow-auto"
          />
        </div>
        
        <div className="space-y-6 overflow-y-auto">
          <div className="card p-6 border-l-4 border-yellow-500">
            <div className="flex items-start gap-3">
              <AlertTriangle className="h-5 w-5 text-yellow-500 shrink-0 mt-0.5" />
              <div>
                <h4 className="text-sm font-bold text-yellow-800 dark:text-yellow-400">Potential Conflict Detected</h4>
                <p className="mt-1 text-xs text-yellow-700 dark:text-yellow-500/80">
                  Rule "deny-ssh-all" (priority 10) is partially shadowed by "allow-ssh-bastion".
                </p>
              </div>
            </div>
          </div>
          
          <div className="card p-6">
            <h3 className="text-sm font-bold mb-4 uppercase tracking-wider text-gray-500">Execution Plan Preview</h3>
            <div className="space-y-3">
              <div className="flex items-center gap-2 text-xs">
                <div className="w-1 h-1 rounded-full bg-green-500"></div>
                <span className="text-green-600 font-medium">+ 2 rules to add</span>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <div className="w-1 h-1 rounded-full bg-red-500"></div>
                <span className="text-red-600 font-medium">- 1 rule to remove</span>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <div className="w-1 h-1 rounded-full bg-yellow-500"></div>
                <span className="text-yellow-600 font-medium">~ 0 rules to modify</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Policies;
