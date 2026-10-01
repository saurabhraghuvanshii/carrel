class Solution {
    public ListNode addDigits(ListNode first, ListNode second) {
        ListNode dummy = new ListNode(0);
        ListNode tail = dummy;
        int carry = 0;
        while (first != null || second != null || carry > 0) {
            int sum = carry;
            if (first != null) {
                sum += first.val;
                first = first.next;
            }
            if (second != null) {
                sum += second.val;
                second = second.next;
            }
            tail.next = new ListNode(sum % 10);
            tail = tail.next;
            carry = sum / 10;
        }
        return dummy.next;
    }
}
