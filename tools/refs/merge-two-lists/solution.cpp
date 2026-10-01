ListNode* mergeLists(ListNode* first, ListNode* second) {
    ListNode dummy(0);
    ListNode* tail = &dummy;
    while (first && second) {
        if (first->val <= second->val) {
            tail->next = first;
            first = first->next;
        } else {
            tail->next = second;
            second = second->next;
        }
        tail = tail->next;
    }
    tail->next = first ? first : second;
    return dummy.next;
}
