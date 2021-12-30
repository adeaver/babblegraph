import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import RadioGroup from '@material-ui/core/RadioGroup';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryRadio } from 'common/components/Radio/Radio';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import { Heading4, Heading6 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';

import {
    ContentNodeType,
    ContentNode,
    ContentNodeBody,
    Heading as HeadingContent,
    Paragraph as ParagraphContent,
    getDefaultContentNodeForType,
} from 'common/api/blog/content.ts';

const styleClasses = makeStyles({
    blogContentEditorContainer: {
        padding: '10px',
    },
});


type BlogContentEditorProps = {
    urlPath: string;
}

const BlogContentEditor = (props: BlogContentEditorProps) => {
    const [ content, setContent ] = useState<Array<ContentNode>>([]);

    const handleRemoveContentAtIndex = (idx: number) => {
        setContent(content.filter((node: ContentNode, nodeIdx: number) => idx !== nodeIdx));
    }
    const handleAppendContentAtIndex = (newNode: ContentNode, idx: number) => {
        setContent(idx === content.length ? (
            content.concat(newNode)
        ) : (
            content.reduce((curr: Array<ContentNode>, nextNode: ContentNode, nodeIdx: number) => (
                nodeIdx === idx ? (
                    curr.concat([newNode, nextNode])
                ) : (
                    curr.concat(nextNode)
                )
            ), [])
        ));
    }
    const handleUpsertContentAtIndex = (newNode: ContentNode, idx: number) => {
        setContent(content.map((nextNode: ContentNode, nodeIdx: number) => (
            nodeIdx === idx ? newNode : nextNode
        )));
    }

    const classes = styleClasses();
    return (
        <Grid container>
            <Grid className={classes.blogContentEditorContainer} item xs={6}>
                <EditorComponent
                    content={content}
                    handleRemoveContentAtIndex={handleRemoveContentAtIndex}
                    handleAppendContentAtIndex={handleAppendContentAtIndex}
                    handleUpsertContentAtIndex={handleUpsertContentAtIndex} />
            </Grid>
            <Grid className={classes.blogContentEditorContainer} item xs={6}>
                <DisplayCard>
                    Hello
                </DisplayCard>
            </Grid>
        </Grid>
    );
}

type EditorComponentProps = {
    content: ContentNode[];

    // TODO: add move content to index
    handleRemoveContentAtIndex: (idx: number) => void;
    handleAppendContentAtIndex: (contentNode: ContentNode, idx: number) => void;
    handleUpsertContentAtIndex: (contentNode: ContentNode, idx: number) => void;
}

const EditorComponent = (props: EditorComponentProps) => {
    const [ newNodeType, setNewNodeType ] = useState<ContentNodeType>(ContentNodeType.Heading);

    const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNewNodeType((event.target as HTMLInputElement).value as ContentNodeType);
    };

    const handleAddNewNode = () => {
        props.handleAppendContentAtIndex(getDefaultContentNodeForType(newNodeType), props.content.length);
    }
    const handleUpdateNode = (idx: number) => {
        return (newNode: ContentNode) => {
            props.handleUpsertContentAtIndex(newNode, idx);
        }
    }

    return (
        <DisplayCard>
            <Heading4 color={TypographyColor.Primary}>
                Add a node
            </Heading4>
            <FormControl component="fieldset">
                <RadioGroup aria-label="add-node-type" name="add-node-type1" value={newNodeType} onChange={handleRadioFormChange}>
                    <Grid container>
                        {
                            Object.values(ContentNodeType).map((option: ContentNodeType, idx: number) => (
                                <EditorComponentAddNodeSelector
                                    key={`add-node-selector-${idx}`}
                                    value={option} />
                            ))
                        }
                    </Grid>
                </RadioGroup>
            </FormControl>
            <PrimaryButton
                type="submit"
                onClick={handleAddNewNode}>
                Submit
            </PrimaryButton>
            {
                props.content.map((node: ContentNode, idx: number) => {
                    if (node.type === ContentNodeType.Heading) {
                        return (
                            <HeadingNodeEditor
                                body={node.body as HeadingContent}
                                handleUpdate={handleUpdateNode(idx)} />
                        );
                    } else if (node.type === ContentNodeType.Paragraph) {
                        return (
                            <ParagraphNodeEditor
                                body={node.body as ParagraphContent}
                                handleUpdate={handleUpdateNode(idx)} />
                        );
                    }
                })
            }
        </DisplayCard>
    );
}

type EditorComponentAddNodeSelectorProps = {
    value: ContentNodeType;
}

const EditorComponentAddNodeSelector  = (props: EditorComponentAddNodeSelectorProps) => {
    return (
        <Grid item xs={6}>
            <FormControlLabel value={props.value} control={<PrimaryRadio />} label={props.value} />
        </Grid>
    )
}

type HeadingNodeEditorProps = {
    body: HeadingContent;
    handleUpdate: (node: ContentNode) => void;
}

const HeadingNodeEditor = (props: HeadingNodeEditorProps) => {
    const [ text, setText ] = useState<string>(props.body.text);
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setText((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        props.handleUpdate({
            type: ContentNodeType.Heading,
            body: {
                text: text,
            }
        })
    }

    return (
        <div>
            <Heading6 color={TypographyColor.Primary}>
                Edit Heading Node
            </Heading6>
            <PrimaryTextField
                value={text}
                label="Heading Content"
                onChange={handleChange}
                variant="outlined" />
            <PrimaryButton
                type="submit"
                onClick={handleSubmit}>
                Submit
            </PrimaryButton>
        </div>
    );
}

type ParagraphNodeEditorProps = {
    body: ParagraphContent;
    handleUpdate: (node: ContentNode) => void;
}

const ParagraphNodeEditor = (props: ParagraphNodeEditorProps) => {
    const [ text, setText ] = useState<string>(props.body.text);
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setText((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        props.handleUpdate({
            type: ContentNodeType.Paragraph,
            body: {
                text: text,
            }
        })
    }

    return (
        <div>
            <Heading6 color={TypographyColor.Primary}>
                Edit Paragraph Node
            </Heading6>
            <PrimaryTextField
                value={text}
                label="Paragraph Content"
                onChange={handleChange}
                variant="outlined" />
            <PrimaryButton
                type="submit"
                onClick={handleSubmit}>
                Submit
            </PrimaryButton>
        </div>
    );
}

export default BlogContentEditor;
