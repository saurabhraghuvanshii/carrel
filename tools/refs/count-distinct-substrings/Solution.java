// Ported from the owner's tries/UniqueSubstrring.java (insert every suffix,
// then countNodes). Fixed: countNodes counts the root too, which stands for
// the empty stretch, so the answer was one too many (8 for "abab"); and the
// root was static, so suffixes from earlier calls stayed in the tree.
class Solution {
    static class Node {
        Node[] children = new Node[26];
        boolean eow;
    }

    private Node root;

    public int countDistinct(String str) {
        root = new Node();
        for (int i = 0; i < str.length(); i++) {
            insert(str.substring(i));
        }
        return countNodes(root) - 1;
    }

    private void insert(String word) {
        Node curr = root;
        for (int level = 0; level < word.length(); level++) {
            int idx = word.charAt(level) - 'a';
            if (curr.children[idx] == null) {
                curr.children[idx] = new Node();
            }
            curr = curr.children[idx];
        }
        curr.eow = true;
    }

    public static int countNodes(Node root) {
        if (root == null) {
            return 0;
        }
        int count = 0;
        for (int i = 0; i < 26; i++) {
            if (root.children[i] != null) {
                count += countNodes(root.children[i]);
            }
        }
        return count + 1;
    }
}
