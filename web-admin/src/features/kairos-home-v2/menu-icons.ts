import {
  Activity,
  Banknote,
  BarChart3,
  Briefcase,
  Building2,
  Calendar,
  Coins,
  FileText,
  Gauge,
  Home,
  Landmark,
  LayoutDashboard,
  LineChart,
  Megaphone,
  Package,
  Percent,
  PieChart,
  PiggyBank,
  Receipt,
  Scale,
  ShoppingCart,
  Target,
  TrendingUp,
  Truck,
  Users,
  Wallet,
} from "lucide-svelte";

/**
 * Icons a canvas can pick for the Home v2 menu with `metadata: kairos_icon:`
 * (lucide names). A closed set on purpose: one icon library and one size, so
 * the menu stays even, and the bundle does not carry all of lucide. An unset
 * or unknown name falls back to the generic dashboard icon.
 */
export const HOME_V2_MENU_ICONS = {
  activity: Activity,
  banknote: Banknote,
  "bar-chart": BarChart3,
  briefcase: Briefcase,
  building: Building2,
  calendar: Calendar,
  coins: Coins,
  "file-text": FileText,
  gauge: Gauge,
  home: Home,
  landmark: Landmark,
  "line-chart": LineChart,
  megaphone: Megaphone,
  package: Package,
  percent: Percent,
  "pie-chart": PieChart,
  "piggy-bank": PiggyBank,
  receipt: Receipt,
  scale: Scale,
  "shopping-cart": ShoppingCart,
  target: Target,
  "trending-up": TrendingUp,
  truck: Truck,
  users: Users,
  wallet: Wallet,
} as const;

export function homeV2MenuIcon(name: string | undefined) {
  return (
    HOME_V2_MENU_ICONS[name as keyof typeof HOME_V2_MENU_ICONS] ??
    LayoutDashboard
  );
}
