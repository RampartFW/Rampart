import React, { useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { 
  LayoutDashboard, 
  ShieldCheck, 
  ListTodo, 
  Zap, 
  Camera, 
  History, 
  Network, 
  Settings,
  Menu,
  X,
  Shield,
  Search,
  Bell
} from 'lucide-react';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const location = useLocation();

  const navigation = [
    { name: 'Dashboard', href: '/', icon: LayoutDashboard },
    { name: 'Policies', href: '/policies', icon: ShieldCheck },
    { name: 'Active Rules', href: '/rules', icon: ListTodo },
    { name: 'Simulator', href: '/simulator', icon: Zap },
    { name: 'Snapshots', href: '/snapshots', icon: Camera },
    { name: 'Audit Log', href: '/audit', icon: History },
    { name: 'Cluster', href: '/cluster', icon: Network },
    { name: 'Settings', href: '/settings', icon: Settings },
  ];

  return (
    <div className="flex h-screen bg-zinc-50 dark:bg-zinc-950 overflow-hidden font-sans antialiased text-zinc-900 dark:text-zinc-100">
      {/* Mobile sidebar overlay */}
      {isSidebarOpen && (
        <div 
          className="fixed inset-0 z-40 bg-zinc-950/60 backdrop-blur-sm lg:hidden" 
          onClick={() => setIsSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside className={cn(
        "fixed inset-y-0 left-0 z-50 w-72 bg-white dark:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0",
        isSidebarOpen ? "translate-x-0" : "-translate-x-full"
      )}>
        <div className="flex items-center gap-3 h-20 px-8 border-b border-zinc-100 dark:border-zinc-800">
          <div className="w-9 h-9 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/20">
            <Shield className="h-5 w-5 text-white" />
          </div>
          <span className="text-xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-zinc-900 to-zinc-500 dark:from-white dark:to-zinc-500">
            RAMPART
          </span>
          <button className="lg:hidden ml-auto p-2 text-zinc-500" onClick={() => setIsSidebarOpen(false)}>
            <X className="h-5 w-5" />
          </button>
        </div>
        
        <nav className="p-6 space-y-2 overflow-y-auto">
          <div className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest px-4 mb-4">Core Management</div>
          {navigation.map((item) => (
            <NavLink
              key={item.name}
              to={item.href}
              className={({ isActive }) => cn(
                "flex items-center px-4 py-3 text-sm font-semibold rounded-xl transition-all duration-200 group",
                isActive 
                  ? "bg-blue-600 text-white shadow-lg shadow-blue-500/25" 
                  : "text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100 hover:bg-zinc-100 dark:hover:bg-white/[0.03]"
              )}
              onClick={() => setIsSidebarOpen(false)}
            >
              <item.icon className={cn(
                "mr-3 h-5 w-5 transition-colors",
                location.pathname === item.href ? "text-white" : "text-zinc-400 group-hover:text-zinc-900 dark:group-hover:text-zinc-100"
              )} />
              {item.name}
            </NavLink>
          ))}
        </nav>

        <div className="absolute bottom-0 left-0 right-0 p-6 border-t border-zinc-100 dark:border-zinc-800">
          <div className="p-4 rounded-2xl bg-zinc-50 dark:bg-white/[0.02] border border-zinc-100 dark:border-zinc-800">
             <div className="flex items-center gap-3">
               <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
               <span className="text-xs font-bold uppercase tracking-wider text-emerald-600 dark:text-emerald-400">Node node-1 (L)</span>
             </div>
             <p className="text-[10px] text-zinc-500 mt-1 font-mono">v0.1.0-production</p>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="flex items-center justify-between h-20 px-8 bg-white/80 dark:bg-zinc-900/80 backdrop-blur-md border-b border-zinc-100 dark:border-zinc-800 relative z-10">
          <button className="lg:hidden p-2 -ml-2 text-zinc-500" onClick={() => setIsSidebarOpen(true)}>
            <Menu className="h-6 w-6" />
          </button>
          
          <div className="flex-1 flex items-center justify-between">
            <div className="hidden md:flex items-center bg-zinc-100 dark:bg-white/[0.05] rounded-xl px-4 py-2 w-96 border border-transparent focus-within:border-blue-500/50 transition-all">
              <Search className="h-4 w-4 text-zinc-400 mr-3" />
              <input 
                type="text" 
                placeholder="Search rules, nodes or events..." 
                className="bg-transparent border-none outline-none text-sm w-full placeholder:text-zinc-500"
              />
            </div>

            <div className="flex items-center gap-4">
               <button className="p-2.5 rounded-xl text-zinc-500 hover:bg-zinc-100 dark:hover:bg-white/[0.05] relative transition-all">
                 <Bell className="h-5 w-5" />
                 <span className="absolute top-2.5 right-2.5 w-2 h-2 bg-rose-500 rounded-full border-2 border-white dark:border-zinc-900" />
               </button>
               <div className="h-8 w-[1px] bg-zinc-100 dark:bg-zinc-800 mx-2" />
               <div className="flex items-center gap-3">
                 <div className="text-right hidden sm:block">
                   <p className="text-sm font-bold tracking-tight">Ersin Koç</p>
                   <p className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">Admin</p>
                 </div>
                 <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-blue-600 border-2 border-white dark:border-zinc-800 shadow-md flex items-center justify-center text-white font-bold">
                   EK
                 </div>
               </div>
            </div>
          </div>
        </header>

        <main className="flex-1 overflow-y-auto p-8 lg:p-12">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};

export default Layout;
