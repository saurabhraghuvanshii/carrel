// Ported from the owner's linkedlist/LinkedList.java (isPalindrome and findMid); head is a parameter
// instead of a static field.
class Solution {
    public boolean isMirror(ListNode head) {
        if (head == null || head.next == null) {
            return true;
        }
        ListNode midNode = findMid(head);
        ListNode prev = null;
        ListNode curr = midNode;
        ListNode next;
        while (curr != null) {
            next = curr.next;
            curr.next = prev;
            prev = curr;
            curr = next;
        }
        ListNode right = prev;
        ListNode left = head;
        while (right != null) {
            if (left.val != right.val) {
                return false;
            }
            left = left.next;
            right = right.next;
        }
        return true;
    }

    private ListNode findMid(ListNode head) {
        ListNode slow = head;
        ListNode fast = head;
        while (fast != null && fast.next != null) {
            slow = slow.next;
            fast = fast.next.next;
        }
        return slow;
    }
}
