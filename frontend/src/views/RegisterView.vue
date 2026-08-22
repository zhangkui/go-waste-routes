<template>
  <div class="auth-card">
    <h1>注册</h1>
    <el-form @submit.prevent="submit">
      <el-form-item><el-input v-model="form.username" placeholder="用户名" /></el-form-item>
      <el-form-item><el-input v-model="form.display_name" placeholder="姓名" /></el-form-item>
      <el-form-item><el-input v-model="form.email" placeholder="邮箱" /></el-form-item>
      <el-form-item><el-input v-model="form.phone" placeholder="电话" /></el-form-item>
      <el-form-item><el-input v-model="form.password" type="password" placeholder="密码" show-password /></el-form-item>
      <el-button type="primary" native-type="submit">注册</el-button>
      <el-button text @click="$router.push('/login')">返回</el-button>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const form = reactive({ username: '', display_name: '', email: '', phone: '', password: '' })

async function submit() {
  await auth.register(form as never)
  router.push('/login')
}
</script>
