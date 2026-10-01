#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iomanip>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <stdexcept>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int v) : val(v), left(nullptr), right(nullptr) {}
};

#include "solution.cpp"

// Reads a tree in level order, "null" for a missing child.
TreeNode* parseTree(const string& line) {
    istringstream in(line);
    vector<string> tok;
    for (string s; in >> s;) tok.push_back(s);
    if (tok.empty() || tok[0] == "null") return nullptr;
    TreeNode* root = new TreeNode(stoi(tok[0]));
    queue<TreeNode*> q;
    q.push(root);
    for (size_t i = 1; i < tok.size();) {
        TreeNode* node = q.front();
        q.pop();
        if (tok[i] != "null") {
            node->left = new TreeNode(stoi(tok[i]));
            q.push(node->left);
        }
        i++;
        if (i < tok.size() && tok[i] != "null") {
            node->right = new TreeNode(stoi(tok[i]));
            q.push(node->right);
        }
        i++;
    }
    return root;
}

size_t countNodes(const string& line) {
    istringstream in(line);
    size_t n = 0;
    for (string s; in >> s;) n += s != "null";
    return n;
}

// Prints a tree in level order as [a, b, null, c]. More than limit nodes
// means the tree has a loop.
string showTree(TreeNode* root, size_t limit) {
    if (root == nullptr) return "[]";
    vector<string> parts;
    queue<TreeNode*> q;
    q.push(root);
    size_t count = 0;
    while (!q.empty()) {
        TreeNode* n = q.front();
        q.pop();
        if (n == nullptr) {
            parts.push_back("null");
            continue;
        }
        if (++count > limit) throw runtime_error("the returned tree has more nodes than it should. Is there a loop?");
        parts.push_back(to_string(n->val));
        q.push(n->left);
        q.push(n->right);
    }
    while (parts.back() == "null") parts.pop_back();
    string out = "[";
    for (size_t i = 0; i < parts.size(); i++) out += (i ? ", " : "") + parts[i];
    return out + "]";
}

// Reads T cases from stdin and prints one line per case. Each case is one tree in level order, "null" for a missing child.
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    string tline;
    getline(cin, tline);
    int t = stoi(tline);
    for (int c = 0; c < t; c++) {
        string line;
        getline(cin, line);
        TreeNode* root = parseTree(line);
        try {
            out << longestPath(root) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
