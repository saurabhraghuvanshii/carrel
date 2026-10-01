void reorder(ListNode* head) {
    ListNode *slow = head, *fast = head->next;
    while (fast && fast->next) {
        slow = slow->next;
        fast = fast->next->next;
    }
    ListNode *curr = slow->next, *prev = nullptr;
    slow->next = nullptr;
    while (curr) {
        ListNode* next = curr->next;
        curr->next = prev;
        prev = curr;
        curr = next;
    }
    for (ListNode *left = head, *right = prev; right;) {
        ListNode *nextL = left->next, *nextR = right->next;
        left->next = right;
        right->next = nextL;
        left = nextL;
        right = nextR;
    }
}
