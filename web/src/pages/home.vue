<template>
    <div class="min-h-screen bg-linear-to-br from-indigo-600 via-purple-600 to-pink-500 text-white">
        <header class="sticky top-0 z-20 backdrop-blur-xl bg-white/10 border-b border-white/20">
            <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
                <div class="flex items-center gap-3">
                    <div class="w-3 h-3 rounded-full bg-emerald-400"></div>
                    <h1 class="text-lg font-semibold tracking-wide">Site Monitor</h1>
                </div>
            </div>
        </header>
        <main class="max-w-7xl mx-auto px-6 py-8 grid grid-cols-1 lg:grid-cols-4 gap-8">
            <section class="lg:col-span-3 space-y-6">
                <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                    <div v-for="item in panelData.projects" :key="item.id"
                        class="rounded-2xl bg-glass backdrop-blur-xl border border-white/20 p-5 shadow-lg hover:shadow-2xl transition">
                        <div class="flex items-start justify-between flex-col">
                            <div class="flex flex-row w-full">
                                <h2 class="font-semibold text-lg">{{ item.name }}</h2>
                                <span class="text-xs px-2 py-1 rounded-full bg-emerald-400/20 h-6 ml-auto"
                                    :class="stateEnums[item.state].color">{{ stateEnums[item.state].text }}</span>
                            </div>
                            <p class="text-xs text-white/70 mt-1">{{ item.repositories }}</p>
                        </div>
                        <div class="mt-4">
                            <div class="flex gap-1">
                                <div class="h-2 flex-1 rounded cursor-pointer" v-for="(history, index) in item.history" :key="index"
                                    :class="historyEnums((history.success_count / history.total_count) * 100)"
                                    @mouseenter="showTooltip($event, history)" @mousemove="moveTooltip($event)"
                                    @mouseleave="hideTooltip" />
                            </div>
                            <p class="text-xs text-white/60 mt-2">最近 24 小时可用率：{{ item.uptimePercentage }}%</p>
                        </div>
                        <div class="mt-4 flex justify-between text-xs text-white/70">
                            <span>检测频率：{{ item.frequency }} 秒</span>
                            <button class="hover:text-white">查看详情 →</button>
                        </div>
                    </div>
                </div>
                <div class="rounded-2xl bg-glass backdrop-blur-xl border border-white/20 p-6 shadow-lg">
                    <h3 class="font-semibold mb-4">最近部署记录</h3>
                    <ul class="space-y-4 text-sm">
                        <li v-for="(log, index) in panelData.logs" :key="index" @click="openLogModal(log.id)"
                            class="cursor-pointer flex flex-col hover:bg-white/10 rounded-lg p-3 transition md:flex-row">
                            <div class="flex flex-col w-auto">
                                <div class="font-medium whitespace-pre-wrap wrap-break-word">{{ log.title }}</div>
                                <div class="text-xs text-white/60 whitespace-pre-wrap wrap-break-word">{{ log.execute }}
                                </div>
                            </div>
                            <span class="text-xs text-white/70 text-left md:ml-auto">{{ log.datetime }}</span>
                        </li>
                    </ul>
                </div>
            </section>
            <aside class="space-y-6">
                <div class="rounded-2xl bg-glass backdrop-blur-xl border border-white/20 p-6 shadow-lg">
                    <h3 class="font-semibold mb-4">整体状态</h3>
                    <div class="space-y-3 text-sm">
                        <div class="flex justify-between"><span>项目总数</span><span>{{ panelData.projectCount }}</span>
                        </div>
                        <div class="flex justify-between"><span>正常运行</span><span class="text-emerald-300">{{
                            panelData.accuracyCount }}</span></div>
                        <div class="flex justify-between"><span>异常</span><span class="text-red-300">{{
                            panelData.errorCount }}</span></div>
                    </div>
                </div>

                <div class="rounded-2xl bg-glass backdrop-blur-xl border border-white/20 p-6 shadow-lg">
                    <h3 class="font-semibold mb-4">Webhook</h3>
                    <code class="block text-xs bg-black/30 rounded-lg p-3 break-all">{{ origin }}/webhook</code>
                </div>
            </aside>
        </main>
        <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm" v-show="logModal">
            <div class="w-full max-w-3xl rounded-2xl bg-black border border-white/20 shadow-2xl">
                <div class="flex items-center justify-between px-6 pt-4 pb-2">
                    <h3 class="font-semibold text-white">
                        {{ detailsData.project }}
                    </h3>
                    <button @click="logModal = false" class="text-white/60 hover:text-white">
                        ✕
                    </button>
                </div>
                <div class="flex items-center gap-2 px-6 pb-3 text-xs font-mono">
                    <button class="px-3 py-1 rounded-md border transition" @click="activeTab = 'meta'"
                        :class="activeTab == 'meta' ? 'bg-emerald-500/20 text-emerald-400 border-emerald-400/40' : 'text-white/50 border-white/10 hover:text-white hover:border-white/30'">基本信息</button>
                    <button class="px-3 py-1 rounded-md border transition" @click="activeTab = 'stdout'"
                        :class="activeTab == 'stdout' ? 'bg-emerald-500/20 text-emerald-400 border-emerald-400/40' : 'text-white/50 border-white/10 hover:text-white hover:border-white/30'">运行输出</button>
                    <button class="px-3 py-1 rounded-md border transition" @click="activeTab = 'stderr'"
                        :class="activeTab == 'stderr' ? 'bg-emerald-500/20 text-emerald-400 border-emerald-400/40' : 'text-white/50 border-white/10 hover:text-white hover:border-white/30'">错误输出</button>
                </div>
                <div v-show="loading"
                    class="mx-6 mb-6 h-80 overflow-auto rounded-lg bg-black/60 p-4 text-xs font-mono items-center justify-center text-white/40 animate-pulse">
                    加载中...
                </div>
                <pre class="mx-6 mb-6 h-80 overflow-auto rounded-lg bg-black/60 p-4 text-xs font-mono text-emerald-400"
                    v-show="activeTab === 'stdout' && !loading">{{ detailsData.stdout }}</pre>
                <pre class="mx-6 mb-6 h-80 overflow-auto rounded-lg bg-black/60 p-4 text-xs font-mono text-red-400"
                    v-show="activeTab === 'stderr' && !loading">{{ detailsData.stderr }}</pre>
                <div class="mx-6 mb-6 h-80 overflow-auto rounded-lg bg-black/60 p-4 text-xs font-mono text-white"
                    v-show="activeTab === 'meta' && !loading">
                    <div class="w-full my-1">[执行命令]: {{ detailsData.command }}</div>
                    <div class="w-full my-1">[开始时间]: {{ detailsData.startTime }}</div>
                    <div class="w-full my-1">[结束时间]: {{ detailsData.endTime }}</div>
                    <div class="w-full h-px my-1 bg-white"></div>
                    <div class="w-full my-1">[执行状态]: <span :class="statusEnums[detailsData.status].color">{{
                        statusEnums[detailsData.status].text }}</span></div>
                    <div class="w-full my-1">[退出码]: {{ detailsData.exitCode }}</div>
                    <div class="w-full my-1">[错误信息]: {{ detailsData.error ? detailsData.error : '暂无错误信息' }}</div>
                </div>
            </div>
        </div>
    </div>
    <div v-if="tooltip.visible" class="fixed z-50 pointer-events-none" ref="tooltipEl" :style="{
        left: tooltip.x + 'px',
        top: tooltip.y + 'px'
    }">
        <div class="rounded-lg bg-black/90 border border-white/20
           px-3 py-2 text-xs text-white shadow-xl">
            <div class="font-mono">
                时间: {{ tooltip.data.hour }}时
            </div>
            <div class="font-mono">
                检查次数: {{ tooltip.data.total_count }}次
            </div>
            <div class="font-mono">
                成功次数: {{ tooltip.data.success_count }}次
            </div>
            <div class="font-mono">
                可用率: {{ Math.round(tooltip.data.success_count / tooltip.data.total_count * 100) }}%
            </div>
            <div class="font-mono">
                平均响应: <span :class="averageResponseColor(tooltip.data.average_response)">{{
                    tooltip.data.average_response }}ms</span>
            </div>
        </div>
    </div>

</template>

<script setup>
import { ref, onMounted } from 'vue';
import httpRequest from '../static/request.js';
import config from '../static/config.js';

const stateEnums = {
    0: { color: "", text: "未监听" },
    1: { color: "text-emerald-300", text: "运行中" },
    2: { color: "text-red-300", text: "异常" }
}

const statusEnums = {
    0: { color: "", text: "待处理" },
    1: { color: "text-emerald-300", text: "运行中" },
    2: { color: "text-emerald-300", text: "成功" },
    3: { color: "text-red-300", text: "失败" },
    4: { color: "text-yellow-300", text: "超时" },
    5: { color: "text-yellow-300", text: "取消" }
}

const historyEnums = (score) => {
    let color = "bg-gray-400"
    if (score > 0) {
        if (score > 90) {
            color = "bg-emerald-400"
        } else if (score > 40) {
            color = "bg-yellow-400"
        } else {
            color = "bg-red-400"
        }
    }
    return color
}

const loading = ref(false)

const panelData = ref({
    projectCount: 0,
    accuracyCount: 0,
    errorCount: 0,
    projects: [],
    logs: []
})

const detailsData = ref({
    id: "",
    project: "",
    command: "",
    status: 0,
    startTime: "",
    endTime: "",
    exitCode: 0,
    error: "",
    stdout: "",
    stderr: ""
})

const origin = window.location.origin

const logModal = ref(false)

const activeTab = ref("meta")

const openLogModal = (id) => {
    loading.value = true
    logModal.value = true
    httpRequest({
        url: config.interface.GetWebhookLogDetails,
        method: 'post',
        data: {
            id: id
        }
    }).then(res => {
        if (res.code == 200 && res.success) {
            detailsData.value = res.data
        } else {
            console.error(res)
        }
    }).catch(err => {
        console.error(err)
    }).finally(() => {
        loading.value = false
    });
}

const tooltipEl = ref(null)
const tooltip = ref({
    visible: false,
    x: 0,
    y: 0,
    data: null,
})

const showTooltip = (e, history) => {
    tooltip.value.visible = true
    tooltip.value.data = history
    moveTooltip(e)
}

const moveTooltip = (e) => {
    const offset = 12
    const el = tooltipEl.value
    if (!el) return

    const { width, height } = el.getBoundingClientRect()

    let x = e.clientX + offset
    let y = e.clientY + offset

    // 右侧溢出，翻到左边
    if (x + width > window.innerWidth) {
        x = e.clientX - width - offset
    }

    // 下方溢出，翻到上面
    if (y + height > window.innerHeight) {
        y = e.clientY - height - offset
    }

    tooltip.value.x = x
    tooltip.value.y = y
}

const hideTooltip = () => {
    tooltip.value.visible = false
}

const averageResponseColor = (average_response) => {
    if (average_response < 100) {
        return "text-emerald-300"
    } else if (average_response < 1000) {
        return "text-yellow-300"
    } else {
        return "text-red-300"
    }
}

onMounted(() => {
    loading.value = true
    httpRequest({
        url: config.interface.GetOverview,
        method: 'post',
        data: {}
    }).then(res => {
        if (res.code == 200 && res.success) {
            panelData.value = res.data
            panelData.value?.projects.forEach(data => {
                let total = 0
                let success = 0
                data.history.forEach(item => {
                    total += item.total_count
                    success += item.success_count
                })
                data.uptimePercentage = parseFloat((success / total) * 100).toFixed(2);
            })
        } else {
            console.error(res)
        }
    }).catch(err => {
        console.error(err)
    }).finally(() => {
        loading.value = false
    });
})
</script>