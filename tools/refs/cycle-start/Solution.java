// Ported from the owner's linkedlist/LinkedList.java (removeCycle): the same meeting-point steps, but it
// returns the node where the loop starts instead of cutting the loop.
class Solution {
    public ListNode loopStart(ListNode head) {
        ListNode slow = head;
        ListNode fast = head;
        boolean cycle = false;
        while (fast != null && fast.next != null) {
            slow = slow.next;
            fast = fast.next.next;
            if (fast == slow) {
                cycle = true;
                break;
            }
        }
        if (!cycle) {
            return null;
        }
        slow = head;
        while (slow != fast) {
            slow = slow.next;
            fast = fast.next;
        }
        return slow;
    }
}
