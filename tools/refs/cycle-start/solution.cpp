ListNode* loopStart(ListNode* head) {
    ListNode *slow = head, *fast = head;
    while (fast && fast->next) {
        slow = slow->next;
        fast = fast->next->next;
        if (slow == fast) {
            for (slow = head; slow != fast; slow = slow->next, fast = fast->next) {}
            return slow;
        }
    }
    return nullptr;
}
