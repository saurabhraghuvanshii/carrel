ListNode* removeFromEnd(ListNode* head, int n) {
    ListNode dummy(0);
    dummy.next = head;
    ListNode *lead = &dummy, *lag = &dummy;
    for (int i = 0; i <= n; i++) lead = lead->next;
    while (lead) {
        lead = lead->next;
        lag = lag->next;
    }
    lag->next = lag->next->next;
    return dummy.next;
}
