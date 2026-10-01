// Ported from the owner's tries/StartsWith.java (insert, search, startsWith)
// as instance methods of the Trie class. Fixed: the root was static, so every
// Trie shared one tree and kept words inserted in earlier tests; it is now a
// field of each Trie.
class Trie {
    static class Node {
        Node[] children = new Node[26];
        boolean eow;
    }

    private Node root = new Node();

    public void insert(String word) {
        int level = 0;
        int len = word.length();
        int idx = 0;
        Node curr = root;
        for (; level < len; level++) {
            idx = word.charAt(level) - 'a';
            if (curr.children[idx] == null) {
                curr.children[idx] = new Node();
            }
            curr = curr.children[idx];
        }
        curr.eow = true;
    }

    public boolean search(String key) {
        int level = 0;
        int len = key.length();
        int idx = 0;
        Node curr = root;
        for (; level < len; level++) {
            idx = key.charAt(level) - 'a';
            if (curr.children[idx] == null) {
                return false;
            }
            curr = curr.children[idx];
        }
        return curr.eow == true;
    }

    public boolean startsWith(String prefix) {
        Node curr = root;
        for (int i = 0; i < prefix.length(); i++) {
            int idx = prefix.charAt(i) - 'a';
            if (curr.children[idx] == null) {
                return false;
            }
            curr = curr.children[idx];
        }
        return true;
    }
}
