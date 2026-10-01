class Solution {
    public ListNode reverseGroups(ListNode head, int k) {
        ListNode dummy = new ListNode(0);
        dummy.next = head;
        ListNode before = dummy;
        while (true) {
            ListNode end = before;
            for (int i = 0; i < k && end != null; i++) {
                end = end.next;
            }
            if (end == null) {
                return dummy.next;
            }
            ListNode first = before.next;
            ListNode after = end.next;
            ListNode prev = after;
            ListNode curr = first;
            while (curr != after) {
                ListNode next = curr.next;
                curr.next = prev;
                prev = curr;
                curr = next;
            }
            before.next = end;
            before = first;
        }
    }
}
