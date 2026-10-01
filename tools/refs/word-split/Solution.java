// Ported from the owner's tries/WordBreak.java. Fixed: the trie root was
// static, so words from earlier calls stayed in it and later texts could split
// with words that are not on their list; it is now made fresh for each call.
// And wordBreak tried every first word again for the same rest of the text,
// which takes exponential time on "aaa...ab"; it now remembers each start.
class Solution {
    static class Node {
        Node children[] = new Node[26];
        boolean eow = false;
    }

    private Node root;
    private Boolean[] memo;

    public boolean canSplitText(String text, String[] words) {
        root = new Node();
        for (String w : words) {
            insert(w);
        }
        memo = new Boolean[text.length() + 1];
        return wordBreak(text, 0);
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

    private boolean search(String key) {
        Node curr = root;
        for (int level = 0; level < key.length(); level++) {
            int idx = key.charAt(level) - 'a';
            if (curr.children[idx] == null) {
                return false;
            }
            curr = curr.children[idx];
        }
        return curr.eow == true;
    }

    private boolean wordBreak(String key, int start) {
        if (start == key.length()) {
            return true;
        }
        if (memo[start] != null) {
            return memo[start];
        }
        boolean found = false;
        for (int i = start + 1; i <= key.length() && !found; i++) {
            found = search(key.substring(start, i)) && wordBreak(key, i);
        }
        return memo[start] = found;
    }
}
