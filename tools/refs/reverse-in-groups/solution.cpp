ListNode* reverseGroups(ListNode* head, int k) {
    ListNode dummy(0);
    dummy.next = head;
    ListNode* before = &dummy;
    while (true) {
        ListNode* end = before;
        for (int i = 0; i < k && end; i++) end = end->next;
        if (!end) return dummy.next;
        ListNode *first = before->next, *after = end->next, *prev = after, *curr = first;
        while (curr != after) {
            ListNode* next = curr->next;
            curr->next = prev;
            prev = curr;
            curr = next;
        }
        before->next = end;
        before = first;
    }
}
