<template>
  <el-container class="shell">
    <el-aside width="240px" class="sidebar">
      <div class="brand">go-waste-routes</div>
      <el-menu router :default-active="$route.path">
        <el-menu-item index="/">仪表盘</el-menu-item>
        <el-menu-item index="/profile">个人资料</el-menu-item>
        <el-menu-item index="/export-center">数据导出</el-menu-item>
        <el-menu-item index="/users">用户管理</el-menu-item>
        <el-menu-item index="/roles">角色管理</el-menu-item>
        <el-menu-item index="/permissions">权限管理</el-menu-item>
        <el-menu-item index="/customers">客户点位</el-menu-item>
        <el-menu-item index="/waste-categories">垃圾类别</el-menu-item>
        <el-menu-item index="/containers">容器管理</el-menu-item>
        <el-menu-item index="/plans">收运计划</el-menu-item>
        <el-menu-item index="/plan-rules">计划规则</el-menu-item>
        <el-menu-item index="/vehicles">车辆管理</el-menu-item>
        <el-menu-item index="/drivers">司机管理</el-menu-item>
        <el-menu-item index="/routes">路线管理</el-menu-item>
        <el-menu-item index="/route-stops">路线站点</el-menu-item>
        <el-menu-item index="/tasks">任务执行</el-menu-item>
        <el-menu-item index="/task-stops">任务站点</el-menu-item>
        <el-menu-item index="/weighing">称重记录</el-menu-item>
        <el-menu-item index="/weighbridges">地磅管理</el-menu-item>
        <el-menu-item index="/abnormalities">异常复核</el-menu-item>
        <el-menu-item index="/invoices">账单管理</el-menu-item>
        <el-menu-item index="/invoice-items">账单明细</el-menu-item>
        <el-menu-item index="/payment-records">收款核销</el-menu-item>
        <el-menu-item index="/audit-logs">审计日志</el-menu-item>
        <el-menu-item index="/system-configs">系统配置</el-menu-item>
        <el-menu-item index="/holiday-configs">节假日配置</el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div>{{ title }}</div>
        <el-button @click="logout">退出</el-button>
      </el-header>
      <el-main class="main"><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const title = computed(() => (route.meta.title as string) ?? 'go-waste-routes')

function logout() {
  auth.logout()
  router.push('/login')
}
</script>
