// Ported from the owner's linkedlist/LinkedList.java (isCycle); head is a parameter instead of a static field.
class Solution {
    public boolean hasLoop(ListNode head) {
        ListNode slow = head;
        ListNode fast = head;
        while (fast != null && fast.next != null) {
            slow = slow.next;
            fast = fast.next.next;
            if (slow == fast) {
                return true;
            }
        }
        return false;
    }
}
