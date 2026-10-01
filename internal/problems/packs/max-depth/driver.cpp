#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
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

// Reads one tree in level order from a line, "null" for a missing child.
TreeNode* parseTree(const string& line) {
    istringstream in(line);
    vector<string> tok;
    for (string s; in >> s;) tok.push_back(s);
    if (tok.empty() || tok[0] == "null") return nullptr;
    TreeNode* root = new TreeNode(stoi(tok[0]));
    queue<TreeNode*> q;
    q.push(root);
    size_t i = 1;
    while (i < tok.size()) {
        TreeNode* node = q.front();
        q.pop();
        if (i < tok.size() && tok[i] != "null") {
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

// Reads T cases from stdin, one tree per line. Prints one line per case.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    string line;
    getline(cin, line);
    int t = stoi(line);
    for (int c = 0; c < t; c++) {
        getline(cin, line);
        TreeNode* root = parseTree(line);
        try {
            out << maxDepth(root) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
