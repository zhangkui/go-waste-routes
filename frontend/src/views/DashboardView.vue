<template>
  <el-row :gutter="16">
    <el-col :span="6"><el-card>今日任务：{{ snapshot.today_tasks }}</el-card></el-col>
    <el-col :span="6"><el-card>待复核：{{ snapshot.pending_abnormality }}</el-card></el-col>
    <el-col :span="6"><el-card>待核销：{{ snapshot.unpaid_invoices }}</el-card></el-col>
    <el-col :span="6"><el-card>总收运量：{{ snapshot.weight_total_ton.toFixed(2) }} 吨</el-card></el-col>
  </el-row>
  <el-card style="margin-top: 16px">
    <template #header>近7日收运趋势</template>
    <el-table :data="snapshot.recent_trend" border>
      <el-table-column prop="date" label="日期" />
      <el-table-column prop="tasks" label="任务数" width="120" />
      <el-table-column prop="weight_ton" label="收运量(吨)" width="160" />
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import http from '@/api/http'

type TrendPoint = {
  date: string
  tasks: number
  weight_ton: number
}

const snapshot = reactive({
  today_tasks: 0,
  claimed_tasks: 0,
  completed_tasks: 0,
  pending_abnormality: 0,
  unpaid_invoices: 0,
  overdue_invoices: 0,
  customer_total: 0,
  route_total: 0,
  weight_total_ton: 0,
  task_status: {} as Record<string, number>,
  invoice_status: {} as Record<string, number>,
  recent_trend: [] as TrendPoint[],
})

onMounted(async () => {
  const response = await http.get('/dashboard')
  Object.assign(snapshot, response.data.data ?? {})
})
</script>
