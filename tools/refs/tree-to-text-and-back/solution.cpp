#include <sstream>
#include <string>
using namespace std;

// Pre-order with "#" for a missing child, read back recursively.
class Codec {
    void write(TreeNode* n, ostringstream& out) {
        if (!n) {
            out << "# ";
            return;
        }
        out << n->val << ' ';
        write(n->left, out);
        write(n->right, out);
    }

    TreeNode* read(istringstream& in) {
        string tok;
        in >> tok;
        if (tok == "#") return nullptr;
        TreeNode* n = new TreeNode(stoi(tok));
        n->left = read(in);
        n->right = read(in);
        return n;
    }

public:
    string serialize(TreeNode* root) {
        ostringstream out;
        write(root, out);
        return out.str();
    }

    TreeNode* deserialize(string text) {
        istringstream in(text);
        return read(in);
    }
};
