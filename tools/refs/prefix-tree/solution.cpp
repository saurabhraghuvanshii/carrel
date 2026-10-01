#include <string>
#include <vector>
using namespace std;

class Trie {
    struct Node {
        int next[26];
        bool end = false;
        Node() { fill(next, next + 26, -1); }
    };
    vector<Node> nodes{Node()};

    int walk(const string& s) {
        int cur = 0;
        for (char c : s) {
            cur = nodes[cur].next[c - 'a'];
            if (cur < 0) return -1;
        }
        return cur;
    }

public:
    void insert(string word) {
        int cur = 0;
        for (char c : word) {
            if (nodes[cur].next[c - 'a'] < 0) {
                nodes[cur].next[c - 'a'] = nodes.size();
                nodes.emplace_back();
            }
            cur = nodes[cur].next[c - 'a'];
        }
        nodes[cur].end = true;
    }

    bool search(string word) {
        int n = walk(word);
        return n >= 0 && nodes[n].end;
    }

    bool startsWith(string prefix) {
        return walk(prefix) >= 0;
    }
};
