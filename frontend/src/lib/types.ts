export interface User {
  id: number;
  username: string;
  full_name: string;
  role: 'admin' | 'warehouse_manager';
  created_at: string;
}

export interface Category {
  id: number;
  name: string;
  created_at: string;
}

export interface Supplier {
  id: number;
  name: string;
  phone: string;
  address: string;
  created_at: string;
}

export interface Product {
  id: number;
  name: string;
  unit: string;
  category_id: number;
  category_name?: string;
  cost_price: number;
  sale_price: number;
  min_stock?: number;
  barcode?: string;
  current_stock?: number;
  stock_value?: number;
  created_at: string;
}

export type TransactionType = 'purchase' | 'return' | 'sale' | 'write_off';

export interface Transaction {
  id: number;
  product_id: number;
  product_name?: string;
  supplier_id?: number | null;
  supplier_name?: string;
  user_id: number;
  user_name?: string;
  type: TransactionType;
  quantity: number;
  unit_price: number;
  total_amount: number;
  note: string;
  expiry_date?: string;
  created_at: string;
}

export interface DashboardAlert {
  type: string; // 'low_stock' | 'expiring_soon'
  product_id: number;
  product_name: string;
  message: string;
  value: string;
}

export interface RecipeIngredient {
  id: number;
  recipe_id: number;
  ingredient_id: number;
  ingredient_name?: string;
  unit?: string;
  quantity: number;
}

export interface Recipe {
  id: number;
  product_id: number;
  product_name?: string;
  name: string;
  ingredients: RecipeIngredient[];
  created_at: string;
}

export type DebtStatus = 'unpaid' | 'partial' | 'paid';

export interface Debt {
  id: number;
  supplier_id: number;
  supplier_name?: string;
  transaction_id: number;
  total_debt: number;
  paid_amount: number;
  status: DebtStatus;
  due_date?: string;
  created_at: string;
}

export interface DashboardSummary {
  total_purchases: number;
  total_returns: number;
  total_sales: number;
  total_write_offs: number;
  total_products: number;
  total_categories: number;
  total_suppliers: number;
  inventory_value: number;
}

export interface ProfitReport {
  total_sales_revenue: number;
  total_cost_of_sold: number;
  total_write_off_loss: number;
  net_profit: number;
  period_from: string;
  period_to: string;
}

export interface TopProduct {
  product_id: number;
  product_name: string;
  total_sold: number;
  total_revenue: number;
}

export interface InventoryItem {
  product_id: number;
  product_name: string;
  unit: string;
  category_name: string;
  cost_price: number;
  sale_price: number;
  current_stock: number;
  stock_value: number;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export const TRANSACTION_LABELS: Record<TransactionType, string> = {
  purchase: 'Kirim',
  return: 'Vozvrat',
  sale: 'Sotuv',
  write_off: 'Spisaniye',
};

export function formatMoney(n: number): string {
  return new Intl.NumberFormat('uz-UZ').format(Math.round(n)) + ' UZS';
}

export function formatDate(d: string): string {
  return new Date(d).toLocaleDateString('uz-UZ', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  });
}

export function formatNumber(n: number): string {
  return new Intl.NumberFormat('uz-UZ', { maximumFractionDigits: 2 }).format(n);
}
