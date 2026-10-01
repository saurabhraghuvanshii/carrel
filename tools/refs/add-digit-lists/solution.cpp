ListNode* addDigits(ListNode* first, ListNode* second) {
    ListNode dummy(0);
    ListNode* tail = &dummy;
    int carry = 0;
    while (first || second || carry) {
        int sum = carry;
        if (first) {
            sum += first->val;
            first = first->next;
        }
        if (second) {
            sum += second->val;
            second = second->next;
        }
        tail->next = new ListNode(sum % 10);
        tail = tail->next;
        carry = sum / 10;
    }
    return dummy.next;
}
