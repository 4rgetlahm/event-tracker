interface Event {
    id: string
    title: string
    description: string
    eventDate: string
    slots: number
}

export interface EventWithUserState extends Event {
    isRegistered: boolean
    slotsLeft: number
}
