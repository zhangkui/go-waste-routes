import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import ExportView from '@/views/ExportView.vue'
import ProfileView from '@/views/ProfileView.vue'
import EntityView from '@/views/EntityView.vue'

const entity = (title: string, apiPath: string) => ({
  component: EntityView,
  meta: { title, apiPath },
})

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: LoginView },
    { path: '/register', component: RegisterView },
    {
      path: '/',
      component: MainLayout,
      children: [
        { path: '', component: DashboardView, meta: { title: '仪表盘' } },
        { path: 'profile', component: ProfileView, meta: { title: '个人资料' } },
        { path: 'export-center', component: ExportView, meta: { title: '数据导出' } },
        { path: 'users', ...entity('用户管理', '/users') },
        { path: 'roles', ...entity('角色管理', '/roles') },
        { path: 'permissions', ...entity('权限管理', '/permissions') },
        { path: 'customers', ...entity('客户点位', '/customers') },
        { path: 'waste-categories', ...entity('垃圾类别', '/waste-categories') },
        { path: 'containers', ...entity('容器管理', '/containers') },
        { path: 'plans', ...entity('收运计划', '/plans') },
        { path: 'plan-rules', ...entity('计划规则', '/plan-rules') },
        { path: 'vehicles', ...entity('车辆管理', '/vehicles') },
        { path: 'drivers', ...entity('司机管理', '/drivers') },
        { path: 'routes', ...entity('路线管理', '/routes') },
        { path: 'route-stops', ...entity('路线站点', '/route-stops') },
        { path: 'tasks', ...entity('任务执行', '/tasks') },
        { path: 'task-stops', ...entity('任务站点', '/task-stops') },
        { path: 'weighing', ...entity('称重记录', '/weighing') },
        { path: 'weighbridges', ...entity('地磅管理', '/weighbridges') },
        { path: 'abnormalities', ...entity('异常复核', '/abnormalities') },
        { path: 'invoices', ...entity('账单管理', '/invoices') },
        { path: 'invoice-items', ...entity('账单明细', '/invoice-items') },
        { path: 'payment-records', ...entity('收款核销', '/payment-records') },
        { path: 'audit-logs', ...entity('审计日志', '/audit-logs') },
        { path: 'system-configs', ...entity('系统配置', '/system-configs') },
        { path: 'holiday-configs', ...entity('节假日配置', '/holiday-configs') },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (!auth.ready && localStorage.getItem('access_token')) {
    await auth.loadMe()
  }
  if (to.path !== '/login' && to.path !== '/register' && !localStorage.getItem('access_token')) {
    return '/login'
  }
})

export default router
