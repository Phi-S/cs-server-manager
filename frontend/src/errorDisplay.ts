import {ref} from "vue";

export var shouldShow = ref(false)
export var title = ref("")
export var messages = ref<string[]>([])

export function Show(titleIn: string, ...message: string[]) {
    title.value = titleIn
    messages.value = message
    shouldShow.value = true
}

export function Hide() {
    title.value = ""
    messages.value = []
    shouldShow.value = false
}